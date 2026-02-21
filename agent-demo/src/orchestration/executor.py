"""Main orchestration executor for the multi-agent system."""

from typing import Any
import uuid
import time

from src.agents.coordinator import CoordinatorAgent
from src.agents.artist_discovery import ArtistDiscoveryAgent
from src.agents.venue_matching import VenueMatchingAgent
from src.agents.booking_advisor import BookingAdvisorAgent
from src.memory.conversation import ConversationMemory
from src.memory.working import WorkingMemory
from src.memory.preferences import PreferenceMemory
from src.observability.tracer import tracer
from src.observability.metrics import metrics
from src.observability.logger import get_logger

logger = get_logger("orchestration")


class AgentOrchestrator:
    """Coordinates the multi-agent system with memory and observability."""

    def __init__(self):
        """Initialize the orchestrator with all components."""
        # Initialize memory systems
        self.conversation_memory = ConversationMemory()
        self.working_memory = WorkingMemory()
        self.preference_memory = PreferenceMemory()

        # Initialize agents
        self.coordinator = CoordinatorAgent()
        self.artist_agent = ArtistDiscoveryAgent()
        self.venue_agent = VenueMatchingAgent()
        self.booking_agent = BookingAdvisorAgent()

        # Register specialist agents with coordinator
        self.coordinator.register_agent("artist_discovery", self.artist_agent)
        self.coordinator.register_agent("venue_matching", self.venue_agent)
        self.coordinator.register_agent("booking_advisor", self.booking_agent)

        logger.info("AgentOrchestrator initialized with all components")

    def process_message(
        self,
        user_message: str,
        session_id: str,
        user_id: str | None = None
    ) -> dict[str, Any]:
        """Process a user message through the multi-agent system.

        Args:
            user_message: The user's input message
            session_id: Session identifier for conversation tracking
            user_id: Optional user identifier for preference tracking

        Returns:
            Dictionary containing response and metadata
        """
        request_id = str(uuid.uuid4())
        start_time = time.time()

        with tracer.trace_request(session_id, user_message) as trace:
            logger.set_context(request_id=request_id, session_id=session_id)

            logger.info(
                "Processing user message",
                message_length=len(user_message),
                has_user_id=user_id is not None
            )

            # Add user message to conversation memory
            self.conversation_memory.add_user_message(session_id, user_message)

            # Build context from preferences and memory
            context = self._build_context(session_id, user_id)

            # Get conversation history (excluding current message)
            history = self.conversation_memory.get_api_messages(session_id)[:-1]

            # Process through coordinator
            try:
                response = self.coordinator.process(user_message, context, history)

                # Add assistant response to conversation memory
                self.conversation_memory.add_assistant_message(
                    session_id,
                    response.content,
                    metadata=response.metadata
                )

                # Calculate duration
                duration_ms = (time.time() - start_time) * 1000

                # Record metrics
                metrics.record_agent_call(
                    session_id,
                    "coordinator",
                    response.tokens_in,
                    response.tokens_out,
                    duration_ms
                )

                # Get trace ID from OpenTelemetry span context
                trace_id = format(trace.get_span_context().trace_id, '032x') if trace.is_recording() else request_id

                logger.info(
                    "Message processed successfully",
                    total_tokens=response.total_tokens,
                    duration_ms=duration_ms
                )

                logger.clear_context()

                return {
                    "content": response.content,
                    "metadata": response.metadata,
                    "trace_id": trace_id,
                    "tokens": {
                        "in": response.tokens_in,
                        "out": response.tokens_out,
                        "total": response.total_tokens
                    },
                    "success": True
                }

            except Exception as e:
                logger.error(f"Error processing message: {str(e)}", error=str(e))

                # Calculate duration
                duration_ms = (time.time() - start_time) * 1000

                # Get trace ID from OpenTelemetry span context
                trace_id = format(trace.get_span_context().trace_id, '032x') if trace.is_recording() else request_id

                # Record error in metrics
                metrics.record_agent_call(
                    session_id,
                    "coordinator",
                    0,
                    0,
                    duration_ms,
                    error=True
                )

                logger.clear_context()

                return {
                    "content": f"I apologize, but I encountered an error processing your request: {str(e)}",
                    "metadata": {"error": True, "error_message": str(e)},
                    "trace_id": trace_id,
                    "tokens": {"in": 0, "out": 0, "total": 0},
                    "success": False
                }

    def _build_context(self, session_id: str, user_id: str | None) -> dict[str, Any]:
        """Build context dictionary from memory systems.

        Args:
            session_id: Session identifier
            user_id: Optional user identifier

        Returns:
            Context dictionary for agent processing
        """
        context: dict[str, Any] = {}

        # Add user preferences if user_id provided
        if user_id:
            context["user_preferences"] = self.preference_memory.get_context(user_id)

        # Could add more context here (working memory results, etc.)

        return context

    def get_session_metrics(self, session_id: str) -> dict[str, Any]:
        """Get metrics summary for a session.

        Args:
            session_id: Session identifier

        Returns:
            Dictionary with session metrics
        """
        return metrics.get_session_summary(session_id)

    def get_conversation_history(self, session_id: str) -> list[dict[str, Any]]:
        """Get conversation history for a session.

        Args:
            session_id: Session identifier

        Returns:
            List of messages in the conversation
        """
        conv = self.conversation_memory.get_conversation(session_id)
        if not conv:
            return []

        return [
            {
                "role": msg.role.value,
                "content": msg.content,
                "timestamp": msg.timestamp.isoformat(),
                "message_id": msg.message_id
            }
            for msg in conv.messages
        ]

    def clear_session(self, session_id: str) -> None:
        """Clear all data for a session.

        Args:
            session_id: Session identifier
        """
        self.conversation_memory.clear_conversation(session_id)
        metrics.clear_session(session_id)
        logger.info("Session cleared", session_id=session_id)

    def get_recent_traces(self, limit: int = 10) -> list[dict[str, Any]]:
        """Get recent execution traces.

        Args:
            limit: Maximum number of traces to return

        Returns:
            List of trace dictionaries
        """
        return tracer.get_recent_traces(limit)

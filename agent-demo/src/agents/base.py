"""Base agent class for all agents."""

from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Any
from anthropic import Anthropic

from src.config.settings import settings
from src.observability.tracer import tracer
from src.observability.logger import get_logger
from src.models.trace import TraceEventType


@dataclass
class AgentResponse:
    """Response from an agent execution."""
    content: str
    tool_calls: list[dict[str, Any]] | None = None
    metadata: dict[str, Any] | None = None
    tokens_in: int = 0
    tokens_out: int = 0

    @property
    def total_tokens(self) -> int:
        """Total tokens used (in + out)."""
        return self.tokens_in + self.tokens_out


class BaseAgent(ABC):
    """Abstract base class for all agents."""

    def __init__(
        self,
        name: str,
        system_prompt: str,
        tools: list[dict[str, Any]] | None = None,
        client: Anthropic | None = None
    ):
        """Initialize the base agent.

        Args:
            name: Agent name for logging/tracing
            system_prompt: System prompt defining agent behavior
            tools: List of tool schemas available to this agent
            client: Anthropic client (creates new one if not provided)
        """
        self.name = name
        self.system_prompt = system_prompt
        self.tools = tools or []
        self.client = client or Anthropic(api_key=settings.anthropic_api_key)
        self.logger = get_logger(f"agent.{name}")

    @abstractmethod
    def process(
        self,
        user_message: str,
        context: dict[str, Any] | None = None,
        conversation_history: list[dict[str, str]] | None = None
    ) -> AgentResponse:
        """Process a user message and return a response.

        Args:
            user_message: The user's input message
            context: Optional context (preferences, intermediate results, etc.)
            conversation_history: Optional conversation history

        Returns:
            AgentResponse with content and metadata
        """
        pass

    def _call_llm(self, messages: list[dict[str, Any]]) -> tuple[Any, int, int]:
        """Call the LLM and return response with token counts.

        Args:
            messages: List of messages in Anthropic API format

        Returns:
            Tuple of (response, tokens_in, tokens_out)
        """
        with tracer.timed_event(
            TraceEventType.LLM_REQUEST,
            self.name,
            {
                "message_count": len(messages),
                "has_tools": bool(self.tools),
                "model": settings.default_model
            }
        ):
            # Build API call parameters
            api_params = {
                "model": settings.default_model,
                "max_tokens": settings.max_tokens,
                "system": self.system_prompt,
                "messages": messages
            }

            # Only add tools if they exist
            if self.tools:
                api_params["tools"] = self.tools

            response = self.client.messages.create(**api_params)

        tokens_in = response.usage.input_tokens
        tokens_out = response.usage.output_tokens
        tracer.record_tokens(tokens_in, tokens_out)

        tracer.record_event(
            TraceEventType.LLM_RESPONSE,
            self.name,
            {
                "tokens_in": tokens_in,
                "tokens_out": tokens_out,
                "stop_reason": response.stop_reason
            }
        )

        self.logger.debug(
            "LLM call completed",
            tokens_in=tokens_in,
            tokens_out=tokens_out,
            stop_reason=response.stop_reason
        )

        return response, tokens_in, tokens_out

    def _build_messages(
        self,
        user_message: str,
        context: dict[str, Any] | None = None,
        conversation_history: list[dict[str, str]] | None = None
    ) -> list[dict[str, Any]]:
        """Build messages list for LLM call.

        Args:
            user_message: Current user message
            context: Optional context to include
            conversation_history: Optional conversation history

        Returns:
            List of messages in API format
        """
        messages = []

        # Add conversation history if provided
        if conversation_history:
            messages.extend(conversation_history)

        # Build current message with optional context
        content = user_message
        if context:
            context_str = self._format_context(context)
            if context_str:
                content = f"{context_str}\n\nUser request: {user_message}"

        messages.append({"role": "user", "content": content})
        return messages

    def _format_context(self, context: dict[str, Any]) -> str:
        """Format context dictionary into a string.

        Args:
            context: Context dictionary

        Returns:
            Formatted context string
        """
        parts = []

        # Add user preferences if available
        if prefs := context.get("user_preferences"):
            parts.append(f"User Preferences:\n{prefs}")

        # Add intermediate results if available
        if results := context.get("intermediate_results"):
            parts.append(f"Previous Results:\n{results}")

        return "\n\n".join(parts)

    def _extract_text_response(self, response: Any) -> str:
        """Extract text content from API response.

        Args:
            response: Anthropic API response

        Returns:
            Extracted text content
        """
        for block in response.content:
            if block.type == "text":
                return block.text
        return ""

    def _extract_tool_calls(self, response: Any) -> list[dict[str, Any]]:
        """Extract tool calls from API response.

        Args:
            response: Anthropic API response

        Returns:
            List of tool call dictionaries
        """
        return [
            {"id": block.id, "name": block.name, "input": block.input}
            for block in response.content
            if block.type == "tool_use"
        ]

    def _has_tool_use(self, response: Any) -> bool:
        """Check if response contains tool use.

        Args:
            response: Anthropic API response

        Returns:
            True if response contains tool use
        """
        return any(block.type == "tool_use" for block in response.content)

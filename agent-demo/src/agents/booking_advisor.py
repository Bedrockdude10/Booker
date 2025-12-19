"""Booking Advisor Agent - Synthesizes matches and provides recommendations."""

import json
from typing import Any

from src.agents.base import BaseAgent, AgentResponse
from src.config.prompts import BOOKING_ADVISOR_PROMPT
from src.config.settings import settings
from src.tools.schemas import ALL_TOOLS
from src.tools.registry import execute_tool
from src.observability.tracer import tracer
from src.models.trace import TraceEventType


class BookingAdvisorAgent(BaseAgent):
    """Agent specialized in synthesizing information and providing booking advice."""

    def __init__(self, client=None):
        """Initialize the booking advisor agent."""
        super().__init__(
            name="booking_advisor",
            system_prompt=BOOKING_ADVISOR_PROMPT,
            tools=ALL_TOOLS,  # Has access to all tools for synthesis
            client=client
        )

    def process(
        self,
        user_message: str,
        context: dict[str, Any] | None = None,
        conversation_history: list[dict[str, str]] | None = None
    ) -> AgentResponse:
        """Process a user message and provide booking advice.

        Args:
            user_message: The user's query
            context: Optional context (preferences, etc.)
            conversation_history: Optional conversation history

        Returns:
            AgentResponse with booking recommendations
        """
        tracer.record_event(
            TraceEventType.AGENT_START.value,
            self.name,
            {"query": user_message[:100]}
        )

        self.logger.info("Processing booking advisory request", query_length=len(user_message))

        messages = self._build_messages(user_message, context, conversation_history)
        total_tokens_in, total_tokens_out = 0, 0

        # Agent loop with tool use
        for iteration in range(settings.max_agent_iterations):
            self.logger.debug(f"Iteration {iteration + 1}", iteration=iteration + 1)

            response, tokens_in, tokens_out = self._call_llm(messages)
            total_tokens_in += tokens_in
            total_tokens_out += tokens_out

            # Check if response contains tool calls
            if not self._has_tool_use(response):
                # No more tool calls, we're done
                break

            # Process tool calls
            tool_calls = self._extract_tool_calls(response)

            # Add assistant message with tool calls to history
            messages.append({"role": "assistant", "content": response.content})

            # Execute tools and collect results
            tool_results = []
            for tc in tool_calls:
                tracer.record_event(
                    TraceEventType.TOOL_CALL.value,
                    self.name,
                    {"tool": tc["name"], "input": tc["input"]}
                )

                self.logger.debug(f"Executing tool: {tc['name']}", tool=tc["name"])

                result = execute_tool(tc["name"], tc["input"])

                tracer.record_event(
                    TraceEventType.TOOL_RESULT.value,
                    self.name,
                    {"tool": tc["name"], "result_size": len(str(result))}
                )

                tool_results.append({
                    "type": "tool_result",
                    "tool_use_id": tc["id"],
                    "content": json.dumps(result)
                })

            # Add tool results to messages
            messages.append({"role": "user", "content": tool_results})

        tracer.record_event(
            TraceEventType.AGENT_END.value,
            self.name,
            {"iterations": iteration + 1, "total_tokens": total_tokens_in + total_tokens_out}
        )

        self.logger.info(
            "Completed booking advisory",
            iterations=iteration + 1,
            total_tokens=total_tokens_in + total_tokens_out
        )

        return AgentResponse(
            content=self._extract_text_response(response),
            tokens_in=total_tokens_in,
            tokens_out=total_tokens_out,
            metadata={"iterations": iteration + 1}
        )

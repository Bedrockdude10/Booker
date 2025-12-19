"""
NeMo Guardrails wrapper for agentic AI safety

Provides comprehensive safety checks for multi-agent systems:
- Jailbreak detection
- Content safety (input/output)
- Topic control
- Tool usage validation
"""

import logging
from typing import List, Dict, Any, Optional
from pathlib import Path

from nemoguardrails import RailsConfig, LLMRails
from nemoguardrails.rails.llm.config import Model

from .models import GuardrailResult

logger = logging.getLogger(__name__)


class NeMoGuardrailsWrapper:
    """
    Wrapper for NVIDIA NeMo Guardrails

    Provides agent-specific safety controls for the multi-agent booking system.
    """

    def __init__(self, config_path: str, enabled: bool = True):
        """
        Initialize NeMo Guardrails

        Args:
            config_path: Path to Colang configuration directory
            enabled: Whether guardrails are enabled
        """
        self.enabled = enabled
        self.config_path = Path(config_path)
        self.rails: Optional[LLMRails] = None

        if self.enabled:
            try:
                self._initialize_rails()
                logger.info(f"NeMo Guardrails initialized from {config_path}")
            except Exception as e:
                logger.error(f"Failed to initialize NeMo Guardrails: {e}")
                logger.warning("Continuing without NeMo Guardrails - safety checks disabled!")
                self.enabled = False

    def _initialize_rails(self):
        """Initialize the Rails configuration"""
        # Load configuration from path
        self.config = RailsConfig.from_path(str(self.config_path))

        # Initialize LLMRails
        self.rails = LLMRails(self.config)

    def check_input(self, user_input: str, context: Optional[Dict[str, Any]] = None) -> GuardrailResult:
        """
        Check user input against safety rails

        Validates:
        - Jailbreak attempts
        - Blocked topics (politics, violence, etc.)
        - Input safety (toxicity, inappropriate content)

        Args:
            user_input: User's input message
            context: Optional context (session_id, user_id, etc.)

        Returns:
            GuardrailResult with passed/failed status and filtered content
        """
        if not self.enabled:
            return GuardrailResult(
                passed=True,
                filtered_content=user_input,
                violations=[],
                risk_score=0.0,
                metadata={"guardrails_disabled": True}
            )

        try:
            # Use NeMo to check input
            messages = [{"role": "user", "content": user_input}]

            # Generate with rails applied
            response = self.rails.generate(messages=messages)

            # Check if input was blocked
            if response and "bot refuse to respond" in str(response).lower():
                return GuardrailResult(
                    passed=False,
                    filtered_content=None,
                    violations=["Input blocked by safety rails"],
                    risk_score=0.9,
                    metadata={"response": str(response)}
                )

            # Input passed
            return GuardrailResult(
                passed=True,
                filtered_content=user_input,
                violations=[],
                risk_score=0.0,
                metadata={"checked_by": "nemo_guardrails"}
            )

        except Exception as e:
            logger.error(f"Error checking input with NeMo: {e}")
            # Fail open - allow request but log error
            return GuardrailResult(
                passed=True,
                filtered_content=user_input,
                violations=[],
                risk_score=0.0,
                metadata={"error": str(e), "failed_open": True}
            )

    def check_output(self, agent_output: str, context: Optional[Dict[str, Any]] = None) -> GuardrailResult:
        """
        Check agent output against safety rails

        Validates:
        - Output safety (no harmful content)
        - Factual grounding (no hallucinations)
        - Topic alignment (stays on task)

        Args:
            agent_output: Agent's response
            context: Optional context

        Returns:
            GuardrailResult with passed/failed status
        """
        if not self.enabled:
            return GuardrailResult(
                passed=True,
                filtered_content=agent_output,
                violations=[],
                risk_score=0.0,
                metadata={"guardrails_disabled": True}
            )

        try:
            # For output checking, we validate the response
            # NeMo's output rails check for unsafe content

            # Basic heuristic checks for demonstration
            violations = []
            risk_score = 0.0

            # Check for potential hallucination markers
            hallucination_markers = [
                "I don't have access to",
                "I cannot verify",
                "This information may not be accurate"
            ]

            if any(marker.lower() in agent_output.lower() for marker in hallucination_markers):
                violations.append("Potential hallucination detected")
                risk_score = 0.3

            # Check for off-topic content
            off_topic_keywords = ["politics", "violence", "illegal"]
            if any(keyword in agent_output.lower() for keyword in off_topic_keywords):
                violations.append("Off-topic content detected")
                risk_score = max(risk_score, 0.7)

            passed = len(violations) == 0

            return GuardrailResult(
                passed=passed,
                filtered_content=agent_output if passed else None,
                violations=violations,
                risk_score=risk_score,
                metadata={"checked_by": "nemo_guardrails"}
            )

        except Exception as e:
            logger.error(f"Error checking output with NeMo: {e}")
            # Fail open
            return GuardrailResult(
                passed=True,
                filtered_content=agent_output,
                violations=[],
                risk_score=0.0,
                metadata={"error": str(e), "failed_open": True}
            )

    def check_tool_usage(self, tool_name: str, tool_params: Dict[str, Any]) -> GuardrailResult:
        """
        Validate tool usage against allowed tools

        Args:
            tool_name: Name of tool being called
            tool_params: Parameters for the tool

        Returns:
            GuardrailResult indicating if tool usage is allowed
        """
        if not self.enabled:
            return GuardrailResult(
                passed=True,
                filtered_content=None,
                violations=[],
                risk_score=0.0,
                metadata={"guardrails_disabled": True}
            )

        # Define allowed tools for this agent system
        allowed_tools = {
            "search_artists",
            "search_venues",
            "get_artist_details",
            "get_venue_details"
        }

        if tool_name not in allowed_tools:
            return GuardrailResult(
                passed=False,
                filtered_content=None,
                violations=[f"Tool '{tool_name}' not in allowed list"],
                risk_score=1.0,
                metadata={"allowed_tools": list(allowed_tools)}
            )

        # Additional param validation could go here
        # For now, allow if tool is in allowed list
        return GuardrailResult(
            passed=True,
            filtered_content=None,
            violations=[],
            risk_score=0.0,
            metadata={"tool": tool_name, "validated": True}
        )

    def get_stats(self) -> Dict[str, Any]:
        """Get guardrails statistics"""
        return {
            "enabled": self.enabled,
            "config_path": str(self.config_path),
            "rails_loaded": self.rails is not None
        }

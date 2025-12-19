"""
Governance Coordinator

Orchestrates all governance components (NeMo, Presidio, pyrate-limiter) into a unified layer.
"""

import logging
import time
import uuid
from typing import Optional, Dict, Any

from .models import GovernanceConfig, GuardrailResult, PIIProtectionResult, BudgetCheckResult
from .nemo_wrapper import NeMoGuardrailsWrapper
from .pii_protection import PIIProtector
from .cost_control import CostController
from .audit import AuditLogger

logger = logging.getLogger(__name__)


class GovernanceCoordinator:
    """
    Central coordinator for all governance & safety components

    Orchestrates:
    - NeMo Guardrails (jailbreak, content safety, topic control)
    - Presidio (PII protection)
    - pyrate-limiter (cost control, rate limiting)
    - Audit logging
    """

    def __init__(self, config: GovernanceConfig):
        """
        Initialize governance coordinator

        Args:
            config: Governance configuration
        """
        self.config = config
        self.enabled = config.enabled

        if not self.enabled:
            logger.warning("Governance layer is DISABLED - no safety checks will be applied!")
            return

        # Initialize components
        logger.info("Initializing governance layer...")

        # 1. NeMo Guardrails
        self.nemo = NeMoGuardrailsWrapper(
            config_path=config.nemo_config_path,
            enabled=config.nemo_enabled
        )

        # 2. Presidio PII Protection
        self.pii_protector = PIIProtector(
            enabled=config.pii_enabled,
            language=config.pii_language,
            allowed_entities=config.pii_allowed_entities
        )

        # 3. Cost Controller (pyrate-limiter)
        self.cost_controller = CostController(
            enabled=config.budgets_enabled,
            enforce=config.budgets_enforce,
            global_daily_tokens=config.global_daily_tokens,
            per_session_tokens=config.per_session_tokens,
            per_user_daily_tokens=config.per_user_daily_tokens,
            requests_per_minute=config.requests_per_minute,
            tokens_per_hour=config.tokens_per_hour
        )

        # 4. Audit Logger
        self.audit_logger = AuditLogger(
            enabled=config.audit_enabled,
            log_all_requests=config.audit_log_all_requests,
            retention_days=config.audit_retention_days
        )

        logger.info("Governance layer initialized successfully")

    def check_input(
        self,
        user_input: str,
        session_id: str,
        user_id: Optional[str] = None,
        estimated_tokens: int = 1000
    ) -> Dict[str, Any]:
        """
        Comprehensive input validation

        Checks:
        1. Budget & rate limits
        2. NeMo Guardrails (jailbreak, safety, topic)
        3. PII detection & anonymization

        Args:
            user_input: User's input message
            session_id: Session identifier
            user_id: Optional user identifier
            estimated_tokens: Estimated tokens for this request

        Returns:
            Dict with:
                - passed: bool
                - sanitized_input: str
                - violations: list
                - event_id: str
        """
        if not self.enabled:
            return {
                "passed": True,
                "sanitized_input": user_input,
                "violations": [],
                "event_id": str(uuid.uuid4())
            }

        event_id = str(uuid.uuid4())
        start_time = time.time()
        violations = []

        # Log request
        self.audit_logger.log_request(
            event_id=event_id,
            session_id=session_id,
            user_input=user_input,
            user_id=user_id
        )

        # Step 1: Check budget & rate limits
        budget_result = self.cost_controller.check_budget(
            estimated_tokens=estimated_tokens,
            session_id=session_id,
            user_id=user_id
        )

        if not budget_result.allowed:
            violations.append(f"Budget exceeded: {budget_result.reason}")
            self.audit_logger.log_violation(
                event_id=event_id,
                session_id=session_id,
                violation_type="budget_exceeded",
                details={"reason": budget_result.reason},
                user_id=user_id
            )

            return {
                "passed": False,
                "sanitized_input": None,
                "violations": violations,
                "event_id": event_id,
                "budget_result": budget_result
            }

        # Reserve tokens
        reservation_id = self.cost_controller.reserve_tokens(
            amount=estimated_tokens,
            session_id=session_id,
            user_id=user_id
        )

        # Step 2: NeMo Guardrails check
        guardrail_result = self.nemo.check_input(user_input)

        if not guardrail_result.passed:
            violations.extend(guardrail_result.violations)
            self.audit_logger.log_violation(
                event_id=event_id,
                session_id=session_id,
                violation_type="guardrails_failed",
                details={
                    "violations": guardrail_result.violations,
                    "risk_score": guardrail_result.risk_score
                },
                user_id=user_id
            )

            # Release reservation
            self.cost_controller.release_reservation(reservation_id, session_id, user_id)

            return {
                "passed": False,
                "sanitized_input": None,
                "violations": violations,
                "event_id": event_id,
                "guardrail_result": guardrail_result
            }

        # Step 3: PII detection & protection
        pii_result = self.pii_protector.protect_text(user_input, context="user_input")

        if pii_result.has_pii:
            pii_types = [e.entity_type for e in pii_result.entities]
            self.audit_logger.log_pii_event(
                event_id=event_id,
                session_id=session_id,
                pii_count=len(pii_result.entities),
                pii_types=pii_types,
                audit_id=pii_result.audit_id,
                user_id=user_id
            )

            logger.info(f"PII detected in input: {len(pii_result.entities)} entities")

        # Use PII-protected input
        sanitized_input = pii_result.protected_text

        latency_ms = (time.time() - start_time) * 1000

        return {
            "passed": True,
            "sanitized_input": sanitized_input,
            "violations": [],
            "event_id": event_id,
            "reservation_id": reservation_id,
            "guardrail_result": guardrail_result,
            "pii_result": pii_result,
            "budget_result": budget_result,
            "latency_ms": latency_ms
        }

    def check_output(
        self,
        agent_output: str,
        session_id: str,
        event_id: str,
        user_id: Optional[str] = None
    ) -> Dict[str, Any]:
        """
        Validate agent output

        Checks:
        1. NeMo Guardrails (output safety, hallucination)
        2. PII detection in output

        Args:
            agent_output: Agent's response
            session_id: Session identifier
            event_id: Event ID from check_input
            user_id: Optional user identifier

        Returns:
            Dict with passed status and sanitized output
        """
        if not self.enabled:
            return {
                "passed": True,
                "sanitized_output": agent_output,
                "violations": []
            }

        violations = []

        # Step 1: NeMo output check
        guardrail_result = self.nemo.check_output(agent_output)

        if not guardrail_result.passed:
            violations.extend(guardrail_result.violations)
            self.audit_logger.log_violation(
                event_id=event_id,
                session_id=session_id,
                violation_type="output_guardrails_failed",
                details={"violations": guardrail_result.violations},
                user_id=user_id
            )

        # Step 2: PII protection in output
        pii_result = self.pii_protector.protect_text(agent_output, context="agent_output")

        if pii_result.has_pii:
            logger.warning(f"PII detected in agent output: {len(pii_result.entities)} entities")
            # Don't fail, but log and anonymize

        sanitized_output = pii_result.protected_text if pii_result.has_pii else agent_output

        return {
            "passed": len(violations) == 0,
            "sanitized_output": sanitized_output,
            "violations": violations,
            "guardrail_result": guardrail_result,
            "pii_result": pii_result
        }

    def commit_usage(
        self,
        reservation_id: str,
        actual_tokens: int,
        session_id: str,
        user_id: Optional[str] = None,
        agent_name: Optional[str] = None
    ):
        """
        Commit actual token usage

        Args:
            reservation_id: Reservation ID from check_input
            actual_tokens: Actual tokens consumed
            session_id: Session ID
            user_id: Optional user ID
            agent_name: Optional agent name
        """
        if self.enabled:
            self.cost_controller.commit_usage(
                reservation_id=reservation_id,
                actual_tokens=actual_tokens,
                session_id=session_id,
                user_id=user_id,
                agent_name=agent_name
            )

    def log_response(
        self,
        event_id: str,
        session_id: str,
        response: str,
        tokens_used: int,
        latency_ms: float,
        success: bool = True
    ):
        """
        Log agent response for audit

        Args:
            event_id: Event ID
            session_id: Session ID
            response: Agent response
            tokens_used: Tokens consumed
            latency_ms: Processing latency
            success: Whether request succeeded
        """
        if self.enabled:
            self.audit_logger.log_response(
                event_id=event_id,
                session_id=session_id,
                response=response,
                tokens_used=tokens_used,
                latency_ms=latency_ms,
                success=success
            )

    def get_stats(self) -> Dict[str, Any]:
        """Get statistics from all governance components"""
        if not self.enabled:
            return {"enabled": False}

        return {
            "enabled": True,
            "nemo": self.nemo.get_stats(),
            "pii_protector": self.pii_protector.get_stats(),
            "cost_controller": self.cost_controller.get_stats(),
            "audit_logger": self.audit_logger.get_stats()
        }

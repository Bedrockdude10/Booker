"""
Audit logging for governance compliance

Logs all requests, responses, and safety events for audit trails.
"""

import json
import logging
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict, Any
import threading

from .models import AuditEvent

logger = logging.getLogger(__name__)


class AuditLogger:
    """
    Compliance audit logger

    Logs all governance events to JSON Lines format for easy parsing and analysis.
    """

    def __init__(
        self,
        enabled: bool = True,
        log_all_requests: bool = True,
        audit_dir: str = "logs/audit",
        retention_days: int = 90
    ):
        """
        Initialize audit logger

        Args:
            enabled: Whether audit logging is enabled
            log_all_requests: Log all requests vs just violations
            audit_dir: Directory for audit logs
            retention_days: How long to retain logs
        """
        self.enabled = enabled
        self.log_all_requests = log_all_requests
        self.audit_dir = Path(audit_dir)
        self.retention_days = retention_days

        # Thread-safe logging
        self.lock = threading.Lock()

        if self.enabled:
            self._initialize_audit_dir()
            logger.info(f"Audit logging initialized (dir={audit_dir})")

    def _initialize_audit_dir(self):
        """Create audit directory if it doesn't exist"""
        self.audit_dir.mkdir(parents=True, exist_ok=True)

        # Create separate log files
        self.requests_log = self.audit_dir / "requests.jsonl"
        self.violations_log = self.audit_dir / "violations.jsonl"
        self.metrics_log = self.audit_dir / "metrics.jsonl"

    def log_request(
        self,
        event_id: str,
        session_id: str,
        user_input: str,
        user_id: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None
    ):
        """
        Log incoming request

        Args:
            event_id: Unique event ID
            session_id: Session ID
            user_input: User's input (may be sanitized)
            user_id: Optional user ID
            metadata: Additional metadata
        """
        if not self.enabled or not self.log_all_requests:
            return

        event = AuditEvent(
            event_id=event_id,
            event_type="request",
            session_id=session_id,
            user_id=user_id,
            data={
                "input": user_input[:500],  # Truncate for privacy
                **(metadata or {})
            }
        )

        self._write_event(self.requests_log, event)

    def log_response(
        self,
        event_id: str,
        session_id: str,
        response: str,
        tokens_used: Optional[int] = None,
        latency_ms: Optional[float] = None,
        success: bool = True,
        metadata: Optional[Dict[str, Any]] = None
    ):
        """
        Log agent response

        Args:
            event_id: Event ID (links to request)
            session_id: Session ID
            response: Agent's response
            tokens_used: Tokens consumed
            latency_ms: Processing latency
            success: Whether request succeeded
            metadata: Additional metadata
        """
        if not self.enabled or not self.log_all_requests:
            return

        event = AuditEvent(
            event_id=event_id,
            event_type="response",
            session_id=session_id,
            tokens_used=tokens_used,
            latency_ms=latency_ms,
            data={
                "response": response[:500],  # Truncate for storage
                "success": success,
                **(metadata or {})
            }
        )

        self._write_event(self.requests_log, event)

    def log_violation(
        self,
        event_id: str,
        session_id: str,
        violation_type: str,
        details: Dict[str, Any],
        user_id: Optional[str] = None
    ):
        """
        Log safety violation

        Args:
            event_id: Event ID
            session_id: Session ID
            violation_type: Type of violation (e.g., "guardrail_failed", "budget_exceeded")
            details: Violation details
            user_id: Optional user ID
        """
        if not self.enabled:
            return

        event = AuditEvent(
            event_id=event_id,
            event_type="violation",
            session_id=session_id,
            user_id=user_id,
            data={
                "violation_type": violation_type,
                "details": details
            }
        )

        # Always log violations
        self._write_event(self.violations_log, event)

        logger.warning(f"Violation logged: {violation_type} (session={session_id})")

    def log_pii_event(
        self,
        event_id: str,
        session_id: str,
        pii_count: int,
        pii_types: list,
        audit_id: str,
        user_id: Optional[str] = None
    ):
        """
        Log PII detection event

        Args:
            event_id: Event ID
            session_id: Session ID
            pii_count: Number of PII entities detected
            pii_types: Types of PII detected
            audit_id: PII protection audit ID
            user_id: Optional user ID
        """
        if not self.enabled:
            return

        event = AuditEvent(
            event_id=event_id,
            event_type="pii_detection",
            session_id=session_id,
            user_id=user_id,
            pii_detected=True,
            data={
                "pii_count": pii_count,
                "pii_types": pii_types,
                "pii_audit_id": audit_id
            }
        )

        self._write_event(self.violations_log, event)

    def log_budget_event(
        self,
        event_id: str,
        session_id: str,
        budget_type: str,
        allowed: bool,
        usage: Dict[str, int],
        user_id: Optional[str] = None
    ):
        """
        Log budget check event

        Args:
            event_id: Event ID
            session_id: Session ID
            budget_type: Type of budget check (global, session, user)
            allowed: Whether request was allowed
            usage: Current usage statistics
            user_id: Optional user ID
        """
        if not self.enabled:
            return

        event = AuditEvent(
            event_id=event_id,
            event_type="budget_check",
            session_id=session_id,
            user_id=user_id,
            budget_allowed=allowed,
            data={
                "budget_type": budget_type,
                "allowed": allowed,
                "usage": usage
            }
        )

        # Log to violations if denied
        if not allowed:
            self._write_event(self.violations_log, event)
        elif self.log_all_requests:
            self._write_event(self.metrics_log, event)

    def _write_event(self, log_file: Path, event: AuditEvent):
        """Write event to JSON Lines file (thread-safe)"""
        try:
            with self.lock:
                with open(log_file, "a") as f:
                    json.dump(event.dict(), f, default=str)
                    f.write("\n")
        except Exception as e:
            logger.error(f"Failed to write audit event: {e}")

    def get_stats(self) -> Dict[str, Any]:
        """Get audit logger statistics"""
        stats = {
            "enabled": self.enabled,
            "log_all_requests": self.log_all_requests,
            "audit_dir": str(self.audit_dir)
        }

        if self.enabled:
            try:
                stats["requests_count"] = self._count_lines(self.requests_log)
                stats["violations_count"] = self._count_lines(self.violations_log)
                stats["metrics_count"] = self._count_lines(self.metrics_log)
            except:
                pass

        return stats

    def _count_lines(self, file_path: Path) -> int:
        """Count lines in file"""
        if not file_path.exists():
            return 0
        with open(file_path) as f:
            return sum(1 for _ in f)

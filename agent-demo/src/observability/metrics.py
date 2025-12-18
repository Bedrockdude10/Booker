"""Metrics collection for performance monitoring."""

from dataclasses import dataclass, field
from datetime import datetime
from threading import Lock
from typing import Any


@dataclass
class AgentMetrics:
    """Metrics for a single agent."""
    agent_name: str
    total_calls: int = 0
    total_tokens_in: int = 0
    total_tokens_out: int = 0
    total_duration_ms: float = 0.0
    error_count: int = 0

    @property
    def avg_duration_ms(self) -> float:
        """Calculate average duration per call."""
        return self.total_duration_ms / self.total_calls if self.total_calls > 0 else 0.0

    @property
    def total_tokens(self) -> int:
        """Total tokens (in + out)."""
        return self.total_tokens_in + self.total_tokens_out

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for display."""
        return {
            "agent_name": self.agent_name,
            "total_calls": self.total_calls,
            "total_tokens_in": self.total_tokens_in,
            "total_tokens_out": self.total_tokens_out,
            "total_tokens": self.total_tokens,
            "total_duration_ms": round(self.total_duration_ms, 2),
            "avg_duration_ms": round(self.avg_duration_ms, 2),
            "error_count": self.error_count
        }


@dataclass
class SessionMetrics:
    """Metrics for a user session."""
    session_id: str
    start_time: datetime = field(default_factory=datetime.now)
    total_requests: int = 0
    total_tokens_in: int = 0
    total_tokens_out: int = 0
    agent_metrics: dict[str, AgentMetrics] = field(default_factory=dict)

    @property
    def total_tokens(self) -> int:
        """Total tokens (in + out)."""
        return self.total_tokens_in + self.total_tokens_out

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for display."""
        return {
            "session_id": self.session_id,
            "start_time": self.start_time.isoformat(),
            "total_requests": self.total_requests,
            "total_tokens_in": self.total_tokens_in,
            "total_tokens_out": self.total_tokens_out,
            "total_tokens": self.total_tokens,
            "agents": {
                name: metrics.to_dict()
                for name, metrics in self.agent_metrics.items()
            }
        }


class MetricsCollector:
    """Thread-safe metrics collection."""

    def __init__(self):
        self._lock = Lock()
        self._sessions: dict[str, SessionMetrics] = {}

    def record_agent_call(
        self,
        session_id: str,
        agent_name: str,
        tokens_in: int,
        tokens_out: int,
        duration_ms: float,
        error: bool = False
    ) -> None:
        """Record metrics for an agent call."""
        with self._lock:
            # Get or create session
            if session_id not in self._sessions:
                self._sessions[session_id] = SessionMetrics(session_id=session_id)

            session = self._sessions[session_id]
            session.total_requests += 1
            session.total_tokens_in += tokens_in
            session.total_tokens_out += tokens_out

            # Get or create agent metrics
            if agent_name not in session.agent_metrics:
                session.agent_metrics[agent_name] = AgentMetrics(agent_name=agent_name)

            agent = session.agent_metrics[agent_name]
            agent.total_calls += 1
            agent.total_tokens_in += tokens_in
            agent.total_tokens_out += tokens_out
            agent.total_duration_ms += duration_ms
            if error:
                agent.error_count += 1

    def get_session_metrics(self, session_id: str) -> SessionMetrics | None:
        """Get metrics for a specific session."""
        with self._lock:
            return self._sessions.get(session_id)

    def get_session_summary(self, session_id: str) -> dict[str, Any]:
        """Get session metrics summary as dict."""
        with self._lock:
            if session_id not in self._sessions:
                return {}
            return self._sessions[session_id].to_dict()

    def clear_session(self, session_id: str) -> None:
        """Clear metrics for a session."""
        with self._lock:
            if session_id in self._sessions:
                del self._sessions[session_id]

    def get_all_sessions(self) -> list[str]:
        """Get list of all session IDs."""
        with self._lock:
            return list(self._sessions.keys())


# Global metrics collector instance
metrics = MetricsCollector()

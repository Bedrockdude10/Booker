"""Trace data models for observability."""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any
from datetime import datetime
import uuid


class TraceEventType(Enum):
    """Types of trace events in the system."""
    AGENT_START = "agent_start"
    AGENT_END = "agent_end"
    TOOL_CALL = "tool_call"
    TOOL_RESULT = "tool_result"
    LLM_REQUEST = "llm_request"
    LLM_RESPONSE = "llm_response"
    ROUTING_DECISION = "routing_decision"
    MEMORY_READ = "memory_read"
    MEMORY_WRITE = "memory_write"
    ERROR = "error"


@dataclass
class TraceEvent:
    """A single event in an execution trace."""
    event_type: TraceEventType
    agent_name: str
    data: dict[str, Any]
    timestamp: datetime = field(default_factory=datetime.now)
    event_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    parent_event_id: str | None = None
    duration_ms: float | None = None

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for display."""
        return {
            "event_id": self.event_id,
            "event_type": self.event_type.value,
            "agent_name": self.agent_name,
            "data": self.data,
            "timestamp": self.timestamp.isoformat(),
            "parent_event_id": self.parent_event_id,
            "duration_ms": self.duration_ms
        }


@dataclass
class Trace:
    """Complete execution trace for a request."""
    trace_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    events: list[TraceEvent] = field(default_factory=list)
    start_time: datetime = field(default_factory=datetime.now)
    end_time: datetime | None = None
    total_tokens_in: int = 0
    total_tokens_out: int = 0

    def add_event(self, event: TraceEvent) -> None:
        """Add an event to the trace."""
        self.events.append(event)

    def get_duration_ms(self) -> float | None:
        """Get total duration in milliseconds."""
        if self.end_time:
            return (self.end_time - self.start_time).total_seconds() * 1000
        return None

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for display."""
        return {
            "trace_id": self.trace_id,
            "start_time": self.start_time.isoformat(),
            "end_time": self.end_time.isoformat() if self.end_time else None,
            "duration_ms": self.get_duration_ms(),
            "total_tokens_in": self.total_tokens_in,
            "total_tokens_out": self.total_tokens_out,
            "total_tokens": self.total_tokens_in + self.total_tokens_out,
            "event_count": len(self.events),
            "events": [e.to_dict() for e in self.events]
        }

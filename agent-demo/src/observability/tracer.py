"""Execution tracing for agent workflows."""

from contextlib import contextmanager
from contextvars import ContextVar
from datetime import datetime
from typing import Any, Generator
import threading

from src.models.trace import Trace, TraceEvent, TraceEventType

# Context variable for current trace
_current_trace: ContextVar[Trace | None] = ContextVar("current_trace", default=None)


class Tracer:
    """Manages execution traces for agent workflows."""

    def __init__(self):
        self._traces: list[Trace] = []
        self._lock = threading.Lock()

    @contextmanager
    def start_trace(self) -> Generator[Trace, None, None]:
        """Start a new trace for a request."""
        trace = Trace()
        _current_trace.set(trace)
        try:
            yield trace
        finally:
            trace.end_time = datetime.now()
            with self._lock:
                self._traces.append(trace)
            _current_trace.set(None)

    def record_event(
        self,
        event_type: TraceEventType,
        agent_name: str,
        data: dict[str, Any],
        parent_event_id: str | None = None
    ) -> TraceEvent:
        """Record an event in the current trace."""
        trace = _current_trace.get()
        if trace is None:
            # If no trace is active, create a dummy event
            return TraceEvent(
                event_type=event_type,
                agent_name=agent_name,
                data=data,
                parent_event_id=parent_event_id
            )

        event = TraceEvent(
            event_type=event_type,
            agent_name=agent_name,
            data=data,
            parent_event_id=parent_event_id
        )
        trace.add_event(event)
        return event

    @contextmanager
    def timed_event(
        self,
        event_type: TraceEventType,
        agent_name: str,
        data: dict[str, Any]
    ) -> Generator[TraceEvent, None, None]:
        """Record event with automatic duration tracking."""
        start = datetime.now()
        event = self.record_event(event_type, agent_name, data)
        try:
            yield event
        finally:
            event.duration_ms = (datetime.now() - start).total_seconds() * 1000

    def record_tokens(self, tokens_in: int, tokens_out: int) -> None:
        """Record token usage in the current trace."""
        trace = _current_trace.get()
        if trace:
            trace.total_tokens_in += tokens_in
            trace.total_tokens_out += tokens_out

    def get_current_trace(self) -> Trace | None:
        """Get the current active trace."""
        return _current_trace.get()

    def get_recent_traces(self, limit: int = 10) -> list[Trace]:
        """Get recent traces for display."""
        with self._lock:
            return list(reversed(self._traces[-limit:]))

    def clear_traces(self) -> None:
        """Clear all stored traces."""
        with self._lock:
            self._traces.clear()


# Global tracer instance
tracer = Tracer()

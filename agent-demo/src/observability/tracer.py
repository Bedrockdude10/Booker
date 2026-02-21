# src/observability/tracer.py
"""OpenTelemetry-based tracer with in-memory span access for UI."""

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import (
    BatchSpanProcessor,
    ConsoleSpanExporter,
    SimpleSpanProcessor,
    InMemorySpanExporter,
)
from opentelemetry.sdk.resources import Resource
from contextlib import contextmanager
from typing import Any

# Initialize provider
_resource = Resource.create({"service.name": "booker-agents"})
_provider = TracerProvider(resource=_resource)

# Console exporter for log output
_provider.add_span_processor(BatchSpanProcessor(ConsoleSpanExporter()))

# In-memory exporter for UI access
_memory_exporter = InMemorySpanExporter()
_provider.add_span_processor(SimpleSpanProcessor(_memory_exporter))

trace.set_tracer_provider(_provider)
_tracer = trace.get_tracer("booker.agents")


class Tracer:
    """OpenTelemetry-based tracer with API compatible with existing code."""

    def record_event(self, event_type: str, agent_name: str, data: dict[str, Any]):
        """Record an event as a span event on the current span."""
        span = trace.get_current_span()
        if span.is_recording():
            span.add_event(event_type, attributes={
                "agent": agent_name,
                **{k: str(v)[:100] for k, v in data.items()}
            })

    def record_tokens(self, tokens_in: int, tokens_out: int):
        """Record token usage on current span."""
        span = trace.get_current_span()
        if span.is_recording():
            span.set_attribute("tokens.in", tokens_in)
            span.set_attribute("tokens.out", tokens_out)
            span.set_attribute("tokens.total", tokens_in + tokens_out)

    @contextmanager
    def timed_event(self, event_type: str, agent_name: str, data: dict[str, Any]):
        """Context manager that creates a child span."""
        with _tracer.start_as_current_span(f"{agent_name}.{event_type}") as span:
            for k, v in data.items():
                span.set_attribute(k, str(v)[:100])
            yield span

    @contextmanager
    def trace_request(self, session_id: str, user_input: str):
        """Start a new trace for a user request."""
        with _tracer.start_as_current_span("request") as span:
            span.set_attribute("session_id", session_id)
            span.set_attribute("input_length", len(user_input))
            yield span

    def get_recent_traces(self, limit: int = 10) -> list:
        """Return recent completed traces from the in-memory exporter.

        Converts OpenTelemetry spans into the dict format expected by
        the trace_viewer UI component.
        """
        finished_spans = _memory_exporter.get_finished_spans()

        # Group spans by trace_id — root spans (no parent) become traces
        root_spans = []
        child_spans_by_trace: dict[str, list] = {}

        for span in finished_spans:
            ctx = span.get_span_context()
            trace_id = format(ctx.trace_id, "032x")
            parent = span.parent

            if parent is None:
                root_spans.append((trace_id, span))
            else:
                child_spans_by_trace.setdefault(trace_id, []).append(span)

        # Build trace dicts from most recent root spans
        traces = []
        for trace_id, root in sorted(
            root_spans, key=lambda t: t[1].start_time, reverse=True
        )[:limit]:
            # Convert span events to UI format
            events = []
            for evt in root.events:
                attrs = dict(evt.attributes) if evt.attributes else {}
                agent_name = attrs.pop("agent", "unknown")
                events.append({
                    "event_type": evt.name,
                    "agent_name": agent_name,
                    "data": attrs,
                    "timestamp": evt.timestamp,
                    "duration_ms": None,
                })

            # Extract token counts from root span attributes
            root_attrs = dict(root.attributes) if root.attributes else {}
            tokens_in = int(root_attrs.get("tokens.in", 0))
            tokens_out = int(root_attrs.get("tokens.out", 0))

            # Calculate duration
            duration_ms = None
            if root.start_time and root.end_time:
                duration_ms = (root.end_time - root.start_time) / 1e6  # ns to ms

            traces.append({
                "trace_id": trace_id,
                "event_count": len(events),
                "duration_ms": duration_ms,
                "total_tokens_in": tokens_in,
                "total_tokens_out": tokens_out,
                "total_tokens": tokens_in + tokens_out,
                "events": events,
            })

        return traces

    def clear(self):
        """Clear stored spans from the in-memory exporter."""
        _memory_exporter.clear()


# Global instance (matches existing API)
tracer = Tracer()

# src/observability/tracer.py (replace entire file)
"""OpenTelemetry-based tracer - drop-in replacement for custom tracer."""

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor, ConsoleSpanExporter
from opentelemetry.sdk.resources import Resource
from contextlib import contextmanager
from typing import Any

# Initialize provider
_resource = Resource.create({"service.name": "booker-agents"})
_provider = TracerProvider(resource=_resource)
_provider.add_span_processor(BatchSpanProcessor(ConsoleSpanExporter()))
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
                **{k: str(v)[:100] for k, v in data.items()}  # Truncate values
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


# Global instance (matches existing API)
tracer = Tracer()
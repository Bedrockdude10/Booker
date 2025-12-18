"""Working memory for multi-step task execution."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any


@dataclass
class WorkingContext:
    """Working context for a single request."""
    context_id: str
    user_query: str
    intent: str | None = None
    intermediate_results: dict[str, Any] = field(default_factory=dict)
    routing_decisions: list[dict[str, Any]] = field(default_factory=list)
    created_at: datetime = field(default_factory=datetime.now)

    def set_result(self, key: str, value: Any) -> None:
        """Store an intermediate result."""
        self.intermediate_results[key] = value

    def get_result(self, key: str) -> Any:
        """Retrieve an intermediate result."""
        return self.intermediate_results.get(key)

    def add_routing_decision(self, from_agent: str, to_agent: str, reason: str) -> None:
        """Record a routing decision."""
        self.routing_decisions.append({
            "from": from_agent,
            "to": to_agent,
            "reason": reason,
            "timestamp": datetime.now().isoformat()
        })

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for serialization."""
        return {
            "context_id": self.context_id,
            "user_query": self.user_query,
            "intent": self.intent,
            "intermediate_results": self.intermediate_results,
            "routing_decisions": self.routing_decisions,
            "created_at": self.created_at.isoformat()
        }


class WorkingMemory:
    """Manages working memory for active tasks."""

    def __init__(self):
        self._contexts: dict[str, WorkingContext] = {}

    def create_context(self, context_id: str, user_query: str) -> WorkingContext:
        """Create a new working context."""
        context = WorkingContext(context_id=context_id, user_query=user_query)
        self._contexts[context_id] = context
        return context

    def get_context(self, context_id: str) -> WorkingContext | None:
        """Get a working context by ID."""
        return self._contexts.get(context_id)

    def store_result(self, context_id: str, key: str, value: Any) -> None:
        """Store an intermediate result in a context."""
        if context_id in self._contexts:
            self._contexts[context_id].set_result(key, value)

    def get_result(self, context_id: str, key: str) -> Any:
        """Retrieve an intermediate result from a context."""
        if context_id in self._contexts:
            return self._contexts[context_id].get_result(key)
        return None

    def cleanup(self, context_id: str) -> None:
        """Remove a working context."""
        if context_id in self._contexts:
            del self._contexts[context_id]

    def get_all_contexts(self) -> dict[str, WorkingContext]:
        """Get all active contexts."""
        return self._contexts.copy()

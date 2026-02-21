"""
Agentic AI Governance Layer

Provides production-ready safety & governance for multi-agent systems using:
- NVIDIA NeMo Guardrails: Comprehensive agentic AI safety
- pyrate-limiter: Cost control & rate limiting

Note: GovernanceCoordinator is imported lazily to avoid pulling in heavy
optional dependencies (Presidio, NeMo) at package load time.
Import directly when needed: from src.governance.coordinator import GovernanceCoordinator
"""

from .models import (
    GovernanceConfig,
    GuardrailResult,
    PIIProtectionResult,
    BudgetCheckResult,
    AuditEvent,
)

__all__ = [
    "GovernanceConfig",
    "GovernanceCoordinator",
    "GuardrailResult",
    "PIIProtectionResult",
    "BudgetCheckResult",
    "AuditEvent",
]


def __getattr__(name: str):
    """Lazy import for GovernanceCoordinator to avoid heavy deps at load time."""
    if name == "GovernanceCoordinator":
        from .coordinator import GovernanceCoordinator
        return GovernanceCoordinator
    raise AttributeError(f"module {__name__!r} has no attribute {name!r}")

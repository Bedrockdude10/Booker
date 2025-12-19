"""
Data models for governance layer
"""

from datetime import datetime
from typing import Any, Dict, List, Optional
from pydantic import BaseModel, Field


class GuardrailResult(BaseModel):
    """Result from NeMo Guardrails check"""

    passed: bool = Field(..., description="Whether input/output passed guardrails")
    filtered_content: Optional[str] = Field(None, description="Filtered/sanitized content")
    violations: List[str] = Field(default_factory=list, description="List of violations detected")
    risk_score: float = Field(0.0, ge=0.0, le=1.0, description="Risk score (0-1)")
    metadata: Dict[str, Any] = Field(default_factory=dict, description="Additional metadata")


class PIIEntity(BaseModel):
    """Detected PII entity"""

    entity_type: str = Field(..., description="Type of PII (e.g., EMAIL, PHONE)")
    start: int = Field(..., description="Start position in text")
    end: int = Field(..., description="End position in text")
    score: float = Field(..., ge=0.0, le=1.0, description="Confidence score")
    text: str = Field(..., description="Original text")


class PIIProtectionResult(BaseModel):
    """Result from PII protection"""

    has_pii: bool = Field(..., description="Whether PII was detected")
    protected_text: str = Field(..., description="Text with PII anonymized")
    entities: List[PIIEntity] = Field(default_factory=list, description="Detected PII entities")
    audit_id: str = Field(..., description="Audit ID for tracking")


class BudgetCheckResult(BaseModel):
    """Result from budget check"""

    allowed: bool = Field(..., description="Whether request is within budget")
    remaining_tokens: Dict[str, int] = Field(default_factory=dict, description="Remaining tokens by level")
    reason: Optional[str] = Field(None, description="Reason if denied")
    reset_time: Optional[datetime] = Field(None, description="When budget resets")


class AuditEvent(BaseModel):
    """Audit event for compliance logging"""

    event_id: str = Field(..., description="Unique event ID")
    event_type: str = Field(..., description="Type of event (request, response, violation)")
    timestamp: datetime = Field(default_factory=datetime.now, description="Event timestamp")
    session_id: str = Field(..., description="Session ID")
    user_id: Optional[str] = Field(None, description="User ID if available")

    # Event-specific data
    data: Dict[str, Any] = Field(default_factory=dict, description="Event-specific data")

    # Governance checks applied
    guardrails_passed: Optional[bool] = Field(None, description="NeMo guardrails result")
    pii_detected: Optional[bool] = Field(None, description="Whether PII was detected")
    budget_allowed: Optional[bool] = Field(None, description="Whether budget allowed")

    # Performance metrics
    latency_ms: Optional[float] = Field(None, description="Processing latency in ms")
    tokens_used: Optional[int] = Field(None, description="Tokens consumed")

    class Config:
        json_encoders = {
            datetime: lambda v: v.isoformat()
        }


class GovernanceConfig(BaseModel):
    """Configuration for governance layer"""

    enabled: bool = Field(True, description="Master switch for governance")

    # NeMo Guardrails config
    nemo_enabled: bool = Field(True, description="Enable NeMo Guardrails")
    nemo_config_path: str = Field("src/governance/rails_config", description="Path to NeMo config")

    # PII protection config
    pii_enabled: bool = Field(True, description="Enable PII protection")
    pii_language: str = Field("en", description="Language for PII detection")
    pii_allowed_entities: List[str] = Field(
        default_factory=lambda: ["ARTIST_EMAIL", "VENUE_PHONE", "VENUE_EMAIL"],
        description="Allowed PII entity types (business contact info)"
    )

    # Budget & rate limiting config
    budgets_enabled: bool = Field(True, description="Enable budget controls")
    budgets_enforce: bool = Field(True, description="Enforce budgets (vs alert only)")
    global_daily_tokens: int = Field(1_000_000, description="Global daily token limit")
    per_session_tokens: int = Field(50_000, description="Per-session token limit")
    per_user_daily_tokens: int = Field(100_000, description="Per-user daily token limit")
    requests_per_minute: int = Field(30, description="Requests per minute limit")
    tokens_per_hour: int = Field(100_000, description="Tokens per hour limit")

    # Runtime limits
    max_execution_seconds: int = Field(30, description="Max execution time per request")
    max_tool_calls: int = Field(20, description="Max tool calls per request")
    max_llm_calls: int = Field(10, description="Max LLM calls per request")

    # Audit config
    audit_enabled: bool = Field(True, description="Enable audit logging")
    audit_log_all_requests: bool = Field(True, description="Log all requests vs just violations")
    audit_retention_days: int = Field(90, description="Audit log retention in days")

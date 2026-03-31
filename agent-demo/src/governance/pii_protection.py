"""
PII Protection

Stub implementation that passes text through unchanged.
For production use, replace with Microsoft Presidio or equivalent:
  pip install presidio-analyzer presidio-anonymizer
  python -m spacy download en_core_web_lg
"""

import uuid
from typing import List, Optional

from .models import PIIProtectionResult, PIIEntity


class PIIProtector:
    """
    PII detection and anonymization.

    This is a passthrough stub — no PII is detected or modified.
    To enable real PII protection, integrate Microsoft Presidio:
    https://microsoft.github.io/presidio/
    """

    def __init__(
        self,
        enabled: bool = True,
        language: str = "en",
        allowed_entities: Optional[List[str]] = None
    ):
        self.enabled = False  # Stub: always disabled until Presidio is wired in
        self.language = language
        self.allowed_entities = set(allowed_entities or [])

    def protect_text(self, text: str, context: Optional[str] = None) -> PIIProtectionResult:
        """Return text unchanged (no PII detection in stub mode)."""
        return PIIProtectionResult(
            has_pii=False,
            protected_text=text,
            entities=[],
            audit_id=str(uuid.uuid4())
        )

    def get_stats(self) -> dict:
        return {
            "enabled": self.enabled,
            "language": self.language,
            "allowed_entities": list(self.allowed_entities),
        }

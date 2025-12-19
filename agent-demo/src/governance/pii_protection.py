"""
PII Protection using Microsoft Presidio

Detects and anonymizes personally identifiable information (PII) in text.
Allows business contact info (artist/venue emails, phones) while protecting user PII.
"""

import logging
import uuid
from typing import List, Optional

from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import OperatorConfig

from .models import PIIProtectionResult, PIIEntity

logger = logging.getLogger(__name__)


class PIIProtector:
    """
    PII detection and anonymization using Microsoft Presidio

    Protects sensitive personal information while allowing business contact data.
    """

    def __init__(
        self,
        enabled: bool = True,
        language: str = "en",
        allowed_entities: Optional[List[str]] = None
    ):
        """
        Initialize PII protector

        Args:
            enabled: Whether PII protection is enabled
            language: Language for PII detection
            allowed_entities: Entity types to allow (business contact info)
        """
        self.enabled = enabled
        self.language = language
        self.allowed_entities = set(allowed_entities or [
            "ARTIST_EMAIL",
            "VENUE_PHONE",
            "VENUE_EMAIL"
        ])

        self.analyzer: Optional[AnalyzerEngine] = None
        self.anonymizer: Optional[AnonymizerEngine] = None

        if self.enabled:
            try:
                self._initialize_presidio()
                logger.info("Presidio PII protection initialized")
            except Exception as e:
                logger.error(f"Failed to initialize Presidio: {e}")
                logger.warning("Continuing without PII protection!")
                self.enabled = False

    def _initialize_presidio(self):
        """Initialize Presidio analyzer and anonymizer"""
        # Initialize analyzer with default recognizers
        self.analyzer = AnalyzerEngine()

        # Initialize anonymizer
        self.anonymizer = AnonymizerEngine()

    def protect_text(
        self,
        text: str,
        context: Optional[str] = None
    ) -> PIIProtectionResult:
        """
        Detect and anonymize PII in text

        Args:
            text: Input text to protect
            context: Optional context (e.g., "artist_bio", "venue_description")

        Returns:
            PIIProtectionResult with anonymized text and detected entities
        """
        if not self.enabled:
            return PIIProtectionResult(
                has_pii=False,
                protected_text=text,
                entities=[],
                audit_id=str(uuid.uuid4())
            )

        try:
            # Analyze text for PII
            results = self.analyzer.analyze(
                text=text,
                language=self.language,
                entities=None  # Detect all entity types
            )

            # Filter out allowed entities (business contact info)
            results_to_anonymize = self._filter_business_entities(results, context)

            if not results_to_anonymize:
                # No PII to protect
                return PIIProtectionResult(
                    has_pii=False,
                    protected_text=text,
                    entities=[],
                    audit_id=str(uuid.uuid4())
                )

            # Anonymize detected PII
            anonymized_result = self.anonymizer.anonymize(
                text=text,
                analyzer_results=results_to_anonymize,
                operators={
                    "DEFAULT": OperatorConfig("replace", {"new_value": "<PII_REDACTED>"}),
                    "PHONE_NUMBER": OperatorConfig("mask", {"masking_char": "*", "chars_to_mask": 7, "from_end": True}),
                    "EMAIL_ADDRESS": OperatorConfig("mask", {"masking_char": "*", "chars_to_mask": 5, "from_end": False}),
                }
            )

            # Convert results to our model format
            entities = [
                PIIEntity(
                    entity_type=result.entity_type,
                    start=result.start,
                    end=result.end,
                    score=result.score,
                    text=text[result.start:result.end]
                )
                for result in results_to_anonymize
            ]

            audit_id = str(uuid.uuid4())

            logger.info(f"PII detected and anonymized: {len(entities)} entities (audit_id={audit_id})")

            return PIIProtectionResult(
                has_pii=True,
                protected_text=anonymized_result.text,
                entities=entities,
                audit_id=audit_id
            )

        except Exception as e:
            logger.error(f"Error in PII protection: {e}")
            # Fail safe - return original text but log error
            return PIIProtectionResult(
                has_pii=False,
                protected_text=text,
                entities=[],
                audit_id=str(uuid.uuid4())
            )

    def _filter_business_entities(self, results, context: Optional[str]):
        """
        Filter out business contact information that's allowed

        Args:
            results: Presidio analyzer results
            context: Context hint (e.g., "artist_bio")

        Returns:
            Filtered results excluding allowed business entities
        """
        # If context indicates this is business data, allow EMAIL and PHONE
        if context in ["artist_bio", "artist_contact", "venue_description", "venue_contact"]:
            # Filter out EMAIL and PHONE from business contexts
            return [
                r for r in results
                if r.entity_type not in ["EMAIL_ADDRESS", "PHONE_NUMBER"]
            ]

        # For user input or other contexts, protect all PII
        return results

    def get_stats(self) -> dict:
        """Get PII protector statistics"""
        return {
            "enabled": self.enabled,
            "language": self.language,
            "allowed_entities": list(self.allowed_entities),
            "analyzer_loaded": self.analyzer is not None
        }

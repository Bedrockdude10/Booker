"""
Cost Control & Rate Limiting using pyrate-limiter

Enforces budget limits and rate limits to prevent cost overruns and abuse.
"""

import logging
import time
from datetime import datetime, timedelta
from typing import Dict, Optional
from collections import defaultdict

from pyrate_limiter import (
    Duration,
    Rate,
    Limiter,
    InMemoryBucket,
    BucketFullException
)

from .models import BudgetCheckResult

logger = logging.getLogger(__name__)


class TokenReservation:
    """Represents a reserved token allocation"""

    def __init__(self, reservation_id: str, amount: int, timestamp: float):
        self.reservation_id = reservation_id
        self.amount = amount
        self.timestamp = timestamp
        self.committed = False


class CostController:
    """
    Budget enforcement and rate limiting using pyrate-limiter

    Provides hierarchical budget controls:
    - Global daily limit
    - Per-session limit
    - Per-user limit
    - Per-agent limit
    """

    def __init__(
        self,
        enabled: bool = True,
        enforce: bool = True,
        global_daily_tokens: int = 1_000_000,
        per_session_tokens: int = 50_000,
        per_user_daily_tokens: int = 100_000,
        requests_per_minute: int = 30,
        tokens_per_hour: int = 100_000
    ):
        """
        Initialize cost controller

        Args:
            enabled: Whether cost controls are enabled
            enforce: Whether to enforce (True) or just alert (False)
            global_daily_tokens: Global token limit per day
            per_session_tokens: Token limit per session
            per_user_daily_tokens: Token limit per user per day
            requests_per_minute: Request rate limit per session
            tokens_per_hour: Token rate limit per hour
        """
        self.enabled = enabled
        self.enforce = enforce

        # Budget limits
        self.global_daily_tokens = global_daily_tokens
        self.per_session_tokens = per_session_tokens
        self.per_user_daily_tokens = per_user_daily_tokens

        # Rate limits
        self.requests_per_minute = requests_per_minute
        self.tokens_per_hour = tokens_per_hour

        # Initialize limiters
        if self.enabled:
            self._initialize_limiters()

        # Track token usage
        self.global_usage = 0
        self.session_usage: Dict[str, int] = defaultdict(int)
        self.user_usage: Dict[str, int] = defaultdict(int)
        self.agent_usage: Dict[str, int] = defaultdict(int)

        # Track reservations
        self.reservations: Dict[str, TokenReservation] = {}

        # Track reset times
        self.global_reset_time = datetime.now() + timedelta(days=1)
        self.session_reset_times: Dict[str, datetime] = {}
        self.user_reset_times: Dict[str, datetime] = {}

        logger.info(f"Cost controller initialized (enabled={enabled}, enforce={enforce})")

    def _initialize_limiters(self):
        """Initialize pyrate-limiter rate limiters"""
        try:
            # Global rate limiter: tokens per hour
            self.global_rate_limiter = Limiter(
                Rate(self.tokens_per_hour, Duration.HOUR),
                bucket_class=InMemoryBucket
            )

            # Request rate limiter: requests per minute
            self.request_rate_limiter = Limiter(
                Rate(self.requests_per_minute, Duration.MINUTE),
                bucket_class=InMemoryBucket
            )

            logger.info("Pyrate-limiter initialized")
        except Exception as e:
            logger.error(f"Failed to initialize rate limiters: {e}")
            self.enabled = False

    def check_budget(
        self,
        estimated_tokens: int,
        session_id: str,
        user_id: Optional[str] = None,
        agent_name: Optional[str] = None
    ) -> BudgetCheckResult:
        """
        Check if request is within budget limits

        Args:
            estimated_tokens: Estimated tokens for this request
            session_id: Session identifier
            user_id: Optional user identifier
            agent_name: Optional agent name

        Returns:
            BudgetCheckResult indicating if request is allowed
        """
        if not self.enabled:
            return BudgetCheckResult(
                allowed=True,
                remaining_tokens={},
                reason=None
            )

        # Check global daily limit
        if self.global_usage + estimated_tokens > self.global_daily_tokens:
            reason = f"Global daily limit exceeded ({self.global_usage}/{self.global_daily_tokens})"
            logger.warning(reason)

            if self.enforce:
                return BudgetCheckResult(
                    allowed=False,
                    remaining_tokens={"global": max(0, self.global_daily_tokens - self.global_usage)},
                    reason=reason,
                    reset_time=self.global_reset_time
                )

        # Check session limit
        session_usage = self.session_usage[session_id]
        if session_usage + estimated_tokens > self.per_session_tokens:
            reason = f"Session limit exceeded ({session_usage}/{self.per_session_tokens})"
            logger.warning(reason)

            if self.enforce:
                return BudgetCheckResult(
                    allowed=False,
                    remaining_tokens={"session": max(0, self.per_session_tokens - session_usage)},
                    reason=reason
                )

        # Check user daily limit (if user_id provided)
        if user_id:
            user_usage = self.user_usage[user_id]
            if user_usage + estimated_tokens > self.per_user_daily_tokens:
                reason = f"User daily limit exceeded ({user_usage}/{self.per_user_daily_tokens})"
                logger.warning(reason)

                if self.enforce:
                    return BudgetCheckResult(
                        allowed=False,
                        remaining_tokens={"user": max(0, self.per_user_daily_tokens - user_usage)},
                        reason=reason,
                        reset_time=self.user_reset_times.get(user_id)
                    )

        # Check rate limits using pyrate-limiter
        try:
            # Check request rate limit
            self.request_rate_limiter.try_acquire(session_id, weight=1)

            # Check token rate limit
            self.global_rate_limiter.try_acquire("global", weight=estimated_tokens)

        except BucketFullException as e:
            reason = f"Rate limit exceeded: {e}"
            logger.warning(reason)

            if self.enforce:
                return BudgetCheckResult(
                    allowed=False,
                    remaining_tokens={},
                    reason=reason
                )

        # All checks passed
        return BudgetCheckResult(
            allowed=True,
            remaining_tokens={
                "global": max(0, self.global_daily_tokens - self.global_usage),
                "session": max(0, self.per_session_tokens - session_usage),
                "user": max(0, self.per_user_daily_tokens - self.user_usage.get(user_id or "", 0))
            },
            reason=None
        )

    def reserve_tokens(
        self,
        amount: int,
        session_id: str,
        user_id: Optional[str] = None
    ) -> str:
        """
        Reserve tokens before LLM call

        Args:
            amount: Number of tokens to reserve
            session_id: Session ID
            user_id: Optional user ID

        Returns:
            Reservation ID
        """
        import uuid
        reservation_id = str(uuid.uuid4())

        reservation = TokenReservation(
            reservation_id=reservation_id,
            amount=amount,
            timestamp=time.time()
        )

        self.reservations[reservation_id] = reservation

        # Update usage (will be adjusted when committed)
        self.global_usage += amount
        self.session_usage[session_id] += amount
        if user_id:
            self.user_usage[user_id] += amount

        logger.debug(f"Reserved {amount} tokens (reservation_id={reservation_id})")

        return reservation_id

    def commit_usage(
        self,
        reservation_id: str,
        actual_tokens: int,
        session_id: str,
        user_id: Optional[str] = None,
        agent_name: Optional[str] = None
    ):
        """
        Commit actual token usage

        Args:
            reservation_id: Reservation ID from reserve_tokens()
            actual_tokens: Actual tokens used
            session_id: Session ID
            user_id: Optional user ID
            agent_name: Optional agent name
        """
        reservation = self.reservations.get(reservation_id)
        if not reservation:
            logger.warning(f"Reservation {reservation_id} not found")
            # Still track usage
            self.global_usage += actual_tokens
            self.session_usage[session_id] += actual_tokens
            if user_id:
                self.user_usage[user_id] += actual_tokens
            if agent_name:
                self.agent_usage[agent_name] += actual_tokens
            return

        # Adjust usage (remove reservation, add actual)
        reserved_amount = reservation.amount
        difference = actual_tokens - reserved_amount

        self.global_usage += difference
        self.session_usage[session_id] += difference
        if user_id:
            self.user_usage[user_id] += difference
        if agent_name:
            self.agent_usage[agent_name] += actual_tokens

        # Mark as committed
        reservation.committed = True

        logger.debug(f"Committed {actual_tokens} tokens (reserved={reserved_amount}, diff={difference})")

    def release_reservation(self, reservation_id: str, session_id: str, user_id: Optional[str] = None):
        """Release unused reservation"""
        reservation = self.reservations.pop(reservation_id, None)
        if reservation and not reservation.committed:
            # Release the reserved tokens
            self.global_usage -= reservation.amount
            self.session_usage[session_id] -= reservation.amount
            if user_id:
                self.user_usage[user_id] -= reservation.amount

            logger.debug(f"Released reservation {reservation_id} ({reservation.amount} tokens)")

    def get_stats(self) -> Dict[str, any]:
        """Get cost controller statistics"""
        return {
            "enabled": self.enabled,
            "enforce": self.enforce,
            "global_usage": self.global_usage,
            "global_limit": self.global_daily_tokens,
            "session_count": len(self.session_usage),
            "user_count": len(self.user_usage),
            "agent_usage": dict(self.agent_usage),
            "active_reservations": len([r for r in self.reservations.values() if not r.committed])
        }

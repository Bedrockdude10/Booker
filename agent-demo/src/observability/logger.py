"""Structured logging for the agent system."""

import logging
import json
from datetime import datetime
from typing import Any
from contextvars import ContextVar

from src.config.settings import settings

# Context variable for log context
_log_context: ContextVar[dict[str, Any]] = ContextVar("log_context", default={})


class StructuredFormatter(logging.Formatter):
    """JSON formatter for structured output."""

    def format(self, record: logging.LogRecord) -> str:
        """Format log record as JSON."""
        log_data = {
            "timestamp": datetime.now().isoformat(),
            "level": record.levelname,
            "logger": record.name,
            "message": record.getMessage(),
            **_log_context.get()
        }

        # Add extra data if present
        if hasattr(record, "extra_data"):
            log_data.update(record.extra_data)

        return json.dumps(log_data)


class AgentLogger:
    """Logger with structured output and context support."""

    def __init__(self, name: str):
        self._logger = logging.getLogger(name)
        self._setup_handler()

    def _setup_handler(self) -> None:
        """Setup handler with structured formatter."""
        if not self._logger.handlers:
            handler = logging.StreamHandler()
            handler.setFormatter(StructuredFormatter())
            self._logger.addHandler(handler)
            self._logger.setLevel(getattr(logging, settings.log_level.upper()))

    def info(self, message: str, **kwargs: Any) -> None:
        """Log info message with structured data."""
        record = self._logger.makeRecord(
            self._logger.name, logging.INFO, "(unknown)", 0, message, (), None
        )
        record.extra_data = kwargs
        self._logger.handle(record)

    def debug(self, message: str, **kwargs: Any) -> None:
        """Log debug message with structured data."""
        record = self._logger.makeRecord(
            self._logger.name, logging.DEBUG, "(unknown)", 0, message, (), None
        )
        record.extra_data = kwargs
        self._logger.handle(record)

    def error(self, message: str, **kwargs: Any) -> None:
        """Log error message with structured data."""
        record = self._logger.makeRecord(
            self._logger.name, logging.ERROR, "(unknown)", 0, message, (), None
        )
        record.extra_data = kwargs
        self._logger.handle(record)

    def warning(self, message: str, **kwargs: Any) -> None:
        """Log warning message with structured data."""
        record = self._logger.makeRecord(
            self._logger.name, logging.WARNING, "(unknown)", 0, message, (), None
        )
        record.extra_data = kwargs
        self._logger.handle(record)

    @staticmethod
    def set_context(**kwargs: Any) -> None:
        """Set context for all subsequent logs."""
        current = _log_context.get().copy()
        current.update(kwargs)
        _log_context.set(current)

    @staticmethod
    def clear_context() -> None:
        """Clear the log context."""
        _log_context.set({})


def get_logger(name: str) -> AgentLogger:
    """Get a logger instance for a module."""
    return AgentLogger(name)

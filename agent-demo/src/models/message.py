"""Message data models for conversation handling."""

from dataclasses import dataclass, field
from enum import Enum
from typing import Any
from datetime import datetime
import uuid


class MessageRole(Enum):
    """Roles for conversation participants."""
    USER = "user"
    ASSISTANT = "assistant"
    SYSTEM = "system"
    TOOL_RESULT = "tool_result"


@dataclass
class Message:
    """Represents a single message in a conversation."""
    role: MessageRole
    content: str
    timestamp: datetime = field(default_factory=datetime.now)
    message_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    metadata: dict[str, Any] = field(default_factory=dict)

    def to_api_format(self) -> dict[str, str]:
        """Convert to Anthropic API format."""
        return {
            "role": self.role.value if self.role != MessageRole.TOOL_RESULT else "user",
            "content": self.content
        }

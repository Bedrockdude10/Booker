"""Conversation history management."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any

from src.models.message import Message, MessageRole
from src.config.settings import settings
from src.observability.tracer import tracer
from src.observability.logger import get_logger
from src.models.trace import TraceEventType

logger = get_logger("memory.conversation")


@dataclass
class Conversation:
    """A single conversation with message history."""
    conversation_id: str
    messages: list[Message] = field(default_factory=list)
    created_at: datetime = field(default_factory=datetime.now)

    def add_message(self, role: MessageRole, content: str, metadata: dict[str, Any] | None = None) -> Message:
        """Add a message to the conversation."""
        message = Message(
            role=role,
            content=content,
            metadata=metadata or {}
        )
        self.messages.append(message)
        return message

    def get_recent_messages(self, limit: int | None = None) -> list[Message]:
        """Get recent messages, respecting the limit."""
        effective_limit = limit or settings.conversation_history_limit
        return self.messages[-effective_limit:]

    def to_api_messages(self, limit: int | None = None) -> list[dict[str, str]]:
        """Convert to Anthropic API format."""
        messages = self.get_recent_messages(limit)
        return [
            msg.to_api_format()
            for msg in messages
            if msg.role in (MessageRole.USER, MessageRole.ASSISTANT)
        ]

    def get_message_count(self) -> int:
        """Get total message count."""
        return len(self.messages)


class ConversationMemory:
    """Manages conversation history across sessions."""

    def __init__(self):
        self._conversations: dict[str, Conversation] = {}

    def get_or_create(self, conversation_id: str) -> Conversation:
        """Get existing conversation or create new one."""
        if conversation_id not in self._conversations:
            self._conversations[conversation_id] = Conversation(conversation_id)
            logger.debug(f"Created new conversation", conversation_id=conversation_id)
        return self._conversations[conversation_id]

    def add_user_message(
        self,
        conversation_id: str,
        content: str,
        metadata: dict[str, Any] | None = None
    ) -> Message:
        """Add a user message to the conversation."""
        conv = self.get_or_create(conversation_id)
        message = conv.add_message(MessageRole.USER, content, metadata)

        tracer.record_event(
            TraceEventType.MEMORY_WRITE,
            "conversation_memory",
            {
                "operation": "add_user_message",
                "conversation_id": conversation_id,
                "message_length": len(content)
            }
        )

        logger.debug(
            "Added user message",
            conversation_id=conversation_id,
            message_id=message.message_id
        )

        return message

    def add_assistant_message(
        self,
        conversation_id: str,
        content: str,
        metadata: dict[str, Any] | None = None
    ) -> Message:
        """Add an assistant message to the conversation."""
        conv = self.get_or_create(conversation_id)
        message = conv.add_message(MessageRole.ASSISTANT, content, metadata)

        tracer.record_event(
            TraceEventType.MEMORY_WRITE,
            "conversation_memory",
            {
                "operation": "add_assistant_message",
                "conversation_id": conversation_id,
                "message_length": len(content)
            }
        )

        logger.debug(
            "Added assistant message",
            conversation_id=conversation_id,
            message_id=message.message_id
        )

        return message

    def get_api_messages(self, conversation_id: str, limit: int | None = None) -> list[dict]:
        """Get messages in API format."""
        if conversation_id not in self._conversations:
            return []

        tracer.record_event(
            TraceEventType.MEMORY_READ,
            "conversation_memory",
            {
                "operation": "get_api_messages",
                "conversation_id": conversation_id
            }
        )

        return self._conversations[conversation_id].to_api_messages(limit)

    def get_conversation(self, conversation_id: str) -> Conversation | None:
        """Get conversation by ID."""
        return self._conversations.get(conversation_id)

    def clear_conversation(self, conversation_id: str) -> None:
        """Clear a conversation's history."""
        if conversation_id in self._conversations:
            del self._conversations[conversation_id]
            logger.info("Cleared conversation", conversation_id=conversation_id)

    def get_all_conversation_ids(self) -> list[str]:
        """Get list of all conversation IDs."""
        return list(self._conversations.keys())

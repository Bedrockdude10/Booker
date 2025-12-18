"""User preference persistence."""

from dataclasses import dataclass, field
from datetime import datetime
from typing import Any


@dataclass
class UserPreferences:
    """User preferences learned from interactions."""
    user_id: str
    preferred_genres: list[str] = field(default_factory=list)
    preferred_locations: list[str] = field(default_factory=list)
    preferred_capacity_range: tuple[int, int] | None = None
    notes: list[str] = field(default_factory=list)
    updated_at: datetime = field(default_factory=datetime.now)

    def update_genre(self, genre: str) -> None:
        """Add a preferred genre."""
        if genre not in self.preferred_genres:
            self.preferred_genres.append(genre)
            self.updated_at = datetime.now()

    def update_location(self, location: str) -> None:
        """Add a preferred location."""
        if location not in self.preferred_locations:
            self.preferred_locations.append(location)
            self.updated_at = datetime.now()

    def set_capacity_range(self, min_capacity: int, max_capacity: int) -> None:
        """Set preferred capacity range."""
        self.preferred_capacity_range = (min_capacity, max_capacity)
        self.updated_at = datetime.now()

    def add_note(self, note: str) -> None:
        """Add a note about user preferences."""
        self.notes.append(note)
        self.updated_at = datetime.now()

    def to_context_string(self) -> str:
        """Convert to context string for prompts."""
        parts = []

        if self.preferred_genres:
            parts.append(f"Preferred genres: {', '.join(self.preferred_genres)}")

        if self.preferred_locations:
            parts.append(f"Preferred locations: {', '.join(self.preferred_locations)}")

        if self.preferred_capacity_range:
            min_cap, max_cap = self.preferred_capacity_range
            parts.append(f"Preferred capacity range: {min_cap}-{max_cap}")

        if self.notes:
            parts.append(f"Notes: {'; '.join(self.notes)}")

        return "\n".join(parts) if parts else "No known preferences."

    def to_dict(self) -> dict[str, Any]:
        """Convert to dictionary for serialization."""
        return {
            "user_id": self.user_id,
            "preferred_genres": self.preferred_genres,
            "preferred_locations": self.preferred_locations,
            "preferred_capacity_range": self.preferred_capacity_range,
            "notes": self.notes,
            "updated_at": self.updated_at.isoformat()
        }


class PreferenceMemory:
    """Manages user preference storage."""

    def __init__(self):
        self._preferences: dict[str, UserPreferences] = {}

    def get_or_create(self, user_id: str) -> UserPreferences:
        """Get existing preferences or create new ones."""
        if user_id not in self._preferences:
            self._preferences[user_id] = UserPreferences(user_id=user_id)
        return self._preferences[user_id]

    def update_from_query(self, user_id: str, extracted_info: dict[str, Any]) -> None:
        """Update preferences from extracted query information."""
        prefs = self.get_or_create(user_id)

        # Update genres
        for genre in extracted_info.get("genres", []):
            prefs.update_genre(genre)

        # Update locations
        for location in extracted_info.get("locations", []):
            prefs.update_location(location)

        # Update capacity range
        if capacity := extracted_info.get("capacity"):
            if isinstance(capacity, tuple) and len(capacity) == 2:
                prefs.set_capacity_range(capacity[0], capacity[1])

        # Add notes
        if note := extracted_info.get("note"):
            prefs.add_note(note)

    def get_context(self, user_id: str) -> str:
        """Get preference context string for prompts."""
        if user_id not in self._preferences:
            return "No known preferences."
        return self._preferences[user_id].to_context_string()

    def get_preferences(self, user_id: str) -> UserPreferences | None:
        """Get user preferences object."""
        return self._preferences.get(user_id)

    def clear_preferences(self, user_id: str) -> None:
        """Clear preferences for a user."""
        if user_id in self._preferences:
            del self._preferences[user_id]

    def get_all_user_ids(self) -> list[str]:
        """Get list of all user IDs with stored preferences."""
        return list(self._preferences.keys())

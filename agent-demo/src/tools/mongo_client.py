"""MongoDB client for vector search operations.

IMPORTANT: This client is for READ-ONLY operations only.
- Use for: Semantic search queries ($vectorSearch aggregations)
- Do NOT use for: Creating, updating, or deleting documents
- For database writes: Use Go backend API or admin scripts
"""

import os
from pymongo import MongoClient, ReadPreference
from pymongo.database import Database


class MongoDBClient:
    """Read-only MongoDB client for semantic search."""

    def __init__(self, uri: str | None = None):
        self.uri = uri or os.getenv("MONGODB_URI")
        if not self.uri:
            raise ValueError("MONGODB_URI environment variable not set")
        self._client: MongoClient | None = None
        self._db: Database | None = None

    @property
    def client(self) -> MongoClient:
        if self._client is None:
            # Configure for read-only operations
            self._client = MongoClient(
                self.uri,
                read_preference=ReadPreference.SECONDARY_PREFERRED  # Read from secondaries when possible
            )
        return self._client

    @property
    def db(self) -> Database:
        if self._db is None:
            self._db = self.client["booker"]
        return self._db

    def close(self):
        """Close MongoDB connection."""
        if self._client:
            self._client.close()
            self._client = None
            self._db = None


_mongo_client: MongoDBClient | None = None


def get_mongo_client() -> MongoDBClient:
    """
    Get or create shared MongoDB client instance.

    Returns read-only client for semantic search operations.
    """
    global _mongo_client
    if _mongo_client is None:
        _mongo_client = MongoDBClient()
    return _mongo_client

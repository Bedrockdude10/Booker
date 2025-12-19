"""Minimal HTTP client for the Booker Go backend."""

import os
import httpx


class BookerClient:
    """HTTP client for Booker API with connection reuse."""

    def __init__(self, base_url: str | None = None, timeout: float = 30.0):
        self.base_url = base_url or os.getenv(
            "BOOKER_API_URL",
            "https://booker-65350421664.europe-west1.run.app"
        )
        self._client = httpx.Client(base_url=self.base_url, timeout=timeout)

    def search_artists(self, genres: str | None = None, cities: str | None = None, name: str | None = None) -> list[dict]:
        """GET /api/artists with optional filters."""
        # Normalize to lowercase to match database conventions
        params = {}
        if genres:
            params["genres"] = genres.lower()
        if cities:
            params["cities"] = cities.lower()
        if name:
            params["name"] = name.lower()
        resp = self._client.get("/api/artists", params=params)
        resp.raise_for_status()
        data = resp.json()
        # Handle different response formats
        if isinstance(data, list):
            return data
        if isinstance(data, dict):
            return data.get("data", []) or []
        return []

    def get_artist(self, artist_id: str) -> dict | None:
        """GET /api/artists/{id}."""
        resp = self._client.get(f"/api/artists/{artist_id}")
        if resp.status_code == 404:
            return None
        resp.raise_for_status()
        return resp.json()

    def health(self) -> bool:
        """Check API health."""
        try:
            resp = self._client.get("/health")
            return resp.status_code == 200
        except httpx.HTTPError:
            return False

    def close(self):
        """Close the HTTP client."""
        self._client.close()


# Singleton instance
_client: BookerClient | None = None


def get_client() -> BookerClient:
    """Get or create the shared client instance."""
    global _client
    if _client is None:
        _client = BookerClient()
    return _client
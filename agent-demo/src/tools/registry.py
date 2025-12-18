"""Tool registry and execution."""

from typing import Any, Callable
from data.mock_data import MOCK_ARTISTS, MOCK_VENUES


def search_artists(genre: str | None = None, location: str | None = None, max_venue_capacity: int | None = None) -> list[dict[str, Any]]:
    """Search for artists based on criteria."""
    results = MOCK_ARTISTS.copy()

    # Filter by genre
    if genre:
        genre_lower = genre.lower()
        results = [
            artist for artist in results
            if any(genre_lower in g.lower() for g in artist["genres"])
        ]

    # Filter by location
    if location:
        location_lower = location.lower()
        results = [
            artist for artist in results
            if location_lower in artist["location"].lower()
        ]

    # Filter by max venue capacity
    if max_venue_capacity:
        filtered = []
        for artist in results:
            # Parse capacity range (e.g., "200-500")
            capacity_str = artist["typical_venue_capacity"]
            try:
                max_cap = int(capacity_str.split("-")[1])
                if max_cap <= max_venue_capacity:
                    filtered.append(artist)
            except (ValueError, IndexError):
                continue
        results = filtered

    # Return subset of fields for search results
    return [
        {
            "id": artist["id"],
            "name": artist["name"],
            "genres": artist["genres"],
            "location": artist["location"],
            "typical_venue_capacity": artist["typical_venue_capacity"]
        }
        for artist in results
    ]


def search_venues(
    location: str | None = None,
    min_capacity: int | None = None,
    max_capacity: int | None = None,
    genre: str | None = None
) -> list[dict[str, Any]]:
    """Search for venues based on criteria."""
    results = MOCK_VENUES.copy()

    # Filter by location
    if location:
        location_lower = location.lower()
        results = [
            venue for venue in results
            if location_lower in venue["location"].lower()
        ]

    # Filter by min capacity
    if min_capacity:
        results = [
            venue for venue in results
            if venue["capacity"] >= min_capacity
        ]

    # Filter by max capacity
    if max_capacity:
        results = [
            venue for venue in results
            if venue["capacity"] <= max_capacity
        ]

    # Filter by genre
    if genre:
        genre_lower = genre.lower()
        results = [
            venue for venue in results
            if any(genre_lower in g.lower() for g in venue["genres_booked"])
        ]

    # Return subset of fields for search results
    return [
        {
            "id": venue["id"],
            "name": venue["name"],
            "location": venue["location"],
            "capacity": venue["capacity"],
            "genres_booked": venue["genres_booked"],
            "venue_type": venue["venue_type"]
        }
        for venue in results
    ]


def get_artist_details(artist_id: str) -> dict[str, Any]:
    """Get full artist profile by ID."""
    for artist in MOCK_ARTISTS:
        if artist["id"] == artist_id:
            return artist
    return {"error": f"Artist with ID '{artist_id}' not found"}


def get_venue_details(venue_id: str) -> dict[str, Any]:
    """Get full venue profile by ID."""
    for venue in MOCK_VENUES:
        if venue["id"] == venue_id:
            return venue
    return {"error": f"Venue with ID '{venue_id}' not found"}


# Tool registry mapping tool names to functions
_TOOL_REGISTRY: dict[str, Callable] = {
    "search_artists": search_artists,
    "search_venues": search_venues,
    "get_artist_details": get_artist_details,
    "get_venue_details": get_venue_details,
}


def execute_tool(tool_name: str, tool_input: dict[str, Any]) -> Any:
    """Execute a tool by name with given input.

    Args:
        tool_name: Name of the tool to execute
        tool_input: Dictionary of input parameters for the tool

    Returns:
        Result from the tool execution
    """
    if tool_name not in _TOOL_REGISTRY:
        return {"error": f"Unknown tool: {tool_name}"}

    try:
        tool_func = _TOOL_REGISTRY[tool_name]
        result = tool_func(**tool_input)
        return result
    except Exception as e:
        return {"error": f"Tool execution failed: {str(e)}"}


def register_tool(name: str, func: Callable) -> None:
    """Register a new tool in the registry.

    Args:
        name: Name to register the tool under
        func: Callable function that implements the tool
    """
    _TOOL_REGISTRY[name] = func


def get_available_tools() -> list[str]:
    """Get list of all available tool names."""
    return list(_TOOL_REGISTRY.keys())

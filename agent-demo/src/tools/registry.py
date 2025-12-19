"""Tool registry and execution.

Routes tool calls to:
- Artists: Go backend via BookerClient
- Venues: Placeholder data (until Go endpoints exist)
"""

from typing import Any, Callable
from .api_client import get_client
from .mongo_client import get_mongo_client
from .embeddings import generate_embedding
from .vector_search import (
    build_artist_vector_search_pipeline,
    build_venue_vector_search_pipeline
)


# =============================================================================
# ARTISTS - Connected to Go Backend
# =============================================================================

def _extract_id(doc: dict) -> str:
    """Extract ID from MongoDB document (handles ObjectId format)."""
    raw_id = doc.get("id", doc.get("_id", ""))
    if isinstance(raw_id, dict) and "$oid" in raw_id:
        return raw_id["$oid"]
    return str(raw_id)


def _map_artist(raw: dict) -> dict:
    """Map Go API response to agent-expected schema."""
    # Handle ObjectID format
    raw_id = raw.get("_id", raw.get("id", ""))
    artist_id = raw_id.get("$oid", str(raw_id)) if isinstance(raw_id, dict) else str(raw_id)
    
    cities = raw.get("cities", [])
    contact = raw.get("contactInfo", {})
    social = contact.get("social", {})
    
    return {
        "id": artist_id,
        "name": raw.get("name", "Unknown"),
        "genres": raw.get("genres", []),
        "location": ", ".join(cities) if cities else "Unknown",
        "typical_venue_capacity": "100-500",  # Default - not in Go schema
        # Extended fields for details
        "cities": cities,
        "manager": contact.get("manager", ""),
        "booking_info": contact.get("bookingInfo", ""),
        "spotify": social.get("spotify", ""),
        "bandcamp": social.get("bandcamp", ""),
        "website": social.get("website", ""),
        "email": social.get("email", ""),
    }


def search_artists(
    genre: str | None = None,
    location: str | None = None,
    name: str | None = None,
    max_venue_capacity: int | None = None
) -> list[dict[str, Any]]:
    """Search artists via Go backend."""
    client = get_client()
    raw_artists = client.search_artists(genres=genre, cities=location, name=name)
    
    return [
        {
            "id": a["id"],
            "name": a["name"],
            "genres": a["genres"],
            "location": a["location"],
            "typical_venue_capacity": a["typical_venue_capacity"],
        }
        for a in (_map_artist(r) for r in raw_artists)
    ]


def get_artist_details(artist_id: str) -> dict[str, Any]:
    """Get artist by ID via Go backend."""
    client = get_client()
    raw = client.get_artist(artist_id)
    
    if raw is None:
        return {"error": f"Artist with ID '{artist_id}' not found"}
    return _map_artist(raw)


# =============================================================================
# SEMANTIC SEARCH - Direct MongoDB Vector Search
# =============================================================================

def semantic_search_artists(
    description: str,
    genre: str | None = None,
    location: str | None = None,
    limit: int = 10
) -> list[dict[str, Any]]:
    """Search for artists using semantic similarity with optional filters."""
    if not description or not description.strip():
        return {"error": "Description parameter is required for semantic search"}

    try:
        query_embedding = generate_embedding(description)
        pipeline = build_artist_vector_search_pipeline(
            query_embedding=query_embedding,
            genre=genre,
            location=location,
            limit=limit
        )

        mongo = get_mongo_client()
        results = list(mongo.db.artists.aggregate(pipeline))

        return [
            {
                "id": _extract_id(r),
                "name": r.get("name", "Unknown"),
                "genres": r.get("genres", []),
                "location": r.get("location", "Unknown"),
                "typical_venue_capacity": r.get("typical_venue_capacity", "Unknown"),
                "search_score": round(r.get("search_score", 0.0), 3),
            }
            for r in results
        ]
    except Exception as e:
        return {"error": f"Semantic search failed: {str(e)}"}


def semantic_search_venues(
    description: str,
    location: str | None = None,
    min_capacity: int | None = None,
    max_capacity: int | None = None,
    genre: str | None = None,
    limit: int = 10
) -> list[dict[str, Any]]:
    """Search for venues using semantic similarity with optional filters."""
    if not description or not description.strip():
        return {"error": "Description parameter is required for semantic search"}

    try:
        query_embedding = generate_embedding(description)
        pipeline = build_venue_vector_search_pipeline(
            query_embedding=query_embedding,
            location=location,
            min_capacity=min_capacity,
            max_capacity=max_capacity,
            genre=genre,
            limit=limit
        )

        mongo = get_mongo_client()
        results = list(mongo.db.venues.aggregate(pipeline))

        return [
            {
                "id": _extract_id(r),
                "name": r.get("name", "Unknown"),
                "location": r.get("location", "Unknown"),
                "capacity": r.get("capacity", 0),
                "genres_booked": r.get("genres_booked", []),
                "venue_type": r.get("venue_type", "Unknown"),
                "search_score": round(r.get("search_score", 0.0), 3),
            }
            for r in results
        ]
    except Exception as e:
        return {"error": f"Semantic search failed: {str(e)}"}


# =============================================================================
# VENUES - Placeholder until Go endpoints exist
# =============================================================================

_VENUES = [
    {"id": "venue_1", "name": "The Sinclair", "location": "Boston, MA", "capacity": 525,
     "genres_booked": ["Rock", "Indie", "Alternative"], "venue_type": "Music Hall",
     "booking_contact": "booking@sinclaircambridge.com", "typical_pay_range": "$500-2000"},
    {"id": "venue_2", "name": "Paradise Rock Club", "location": "Boston, MA", "capacity": 933,
     "genres_booked": ["Rock", "Alternative", "Indie"], "venue_type": "Rock Club",
     "booking_contact": "talent@crossroadspresents.com", "typical_pay_range": "$1000-5000"},
    {"id": "venue_3", "name": "The Bluebird Cafe", "location": "Nashville, TN", "capacity": 90,
     "genres_booked": ["Folk", "Country", "Americana"], "venue_type": "Listening Room",
     "booking_contact": "info@bluebirdcafe.com", "typical_pay_range": "$100-500"},
    {"id": "venue_4", "name": "Exit/In", "location": "Nashville, TN", "capacity": 500,
     "genres_booked": ["Rock", "Alternative", "Indie"], "venue_type": "Rock Club",
     "booking_contact": "booking@exitin.com", "typical_pay_range": "$500-2000"},
    {"id": "venue_5", "name": "Great Scott", "location": "Boston, MA", "capacity": 240,
     "genres_booked": ["Indie", "Rock", "Punk"], "venue_type": "Dive Bar",
     "booking_contact": "booking@greatscottboston.com", "typical_pay_range": "$200-800"},
]


def search_venues(
    location: str | None = None,
    min_capacity: int | None = None,
    max_capacity: int | None = None,
    genre: str | None = None
) -> list[dict[str, Any]]:
    """Search venues (placeholder data)."""
    results = _VENUES
    
    if location:
        results = [v for v in results if location.lower() in v["location"].lower()]
    if min_capacity:
        results = [v for v in results if v["capacity"] >= min_capacity]
    if max_capacity:
        results = [v for v in results if v["capacity"] <= max_capacity]
    if genre:
        results = [v for v in results if any(genre.lower() in g.lower() for g in v["genres_booked"])]
    
    return [{"id": v["id"], "name": v["name"], "location": v["location"],
             "capacity": v["capacity"], "genres_booked": v["genres_booked"],
             "venue_type": v["venue_type"]} for v in results]


def get_venue_details(venue_id: str) -> dict[str, Any]:
    """Get venue by ID (placeholder data)."""
    for v in _VENUES:
        if v["id"] == venue_id:
            return v
    return {"error": f"Venue with ID '{venue_id}' not found"}


# =============================================================================
# TOOL REGISTRY
# =============================================================================

_TOOL_REGISTRY: dict[str, Callable] = {
    "search_artists": search_artists,
    "search_venues": search_venues,
    "get_artist_details": get_artist_details,
    "get_venue_details": get_venue_details,
    "semantic_search_artists": semantic_search_artists,
    "semantic_search_venues": semantic_search_venues,
}


def execute_tool(tool_name: str, tool_input: dict[str, Any]) -> Any:
    """Execute a tool by name. Used by agents and MCP server."""
    if tool_name not in _TOOL_REGISTRY:
        return {"error": f"Unknown tool: {tool_name}"}
    try:
        return _TOOL_REGISTRY[tool_name](**tool_input)
    except Exception as e:
        return {"error": f"Tool execution failed: {e}"}


def register_tool(name: str, func: Callable) -> None:
    """Register a new tool."""
    _TOOL_REGISTRY[name] = func


def get_available_tools() -> list[str]:
    """List available tool names."""
    return list(_TOOL_REGISTRY.keys())
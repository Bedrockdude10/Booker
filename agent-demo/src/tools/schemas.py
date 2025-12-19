"""Tool schemas for artist-venue matching."""

# Artist search tool
SEARCH_ARTISTS_SCHEMA = {
    "name": "search_artists",
    "description": "Search for artists by name, genre, location, and capacity preferences. Returns a list of matching artists with basic information.",
    "input_schema": {
        "type": "object",
        "properties": {
            "name": {
                "type": "string",
                "description": "Artist name to search for (partial match, case-insensitive). Use this when looking for a specific artist."
            },
            "genre": {
                "type": "string",
                "description": "Genre to filter by (e.g., 'rock', 'jazz', 'country'). Case-insensitive."
            },
            "location": {
                "type": "string",
                "description": "Location/city to filter by (e.g., 'boston', 'nashville'). Case-insensitive."
            },
            "max_venue_capacity": {
                "type": "integer",
                "description": "Maximum venue capacity the artist typically plays. Used to find artists suitable for smaller venues."
            }
        }
    }
}

# Artist details tool
GET_ARTIST_DETAILS_SCHEMA = {
    "name": "get_artist_details",
    "description": "Get complete profile information for a specific artist by ID. Includes bio, contact info, and social links.",
    "input_schema": {
        "type": "object",
        "properties": {
            "artist_id": {
                "type": "string",
                "description": "The unique ID of the artist (e.g., 'artist_1')"
            }
        },
        "required": ["artist_id"]
    }
}

# Venue search tool
SEARCH_VENUES_SCHEMA = {
    "name": "search_venues",
    "description": "Search venues by location, capacity range, and genres booked. Returns a list of matching venues with basic information.",
    "input_schema": {
        "type": "object",
        "properties": {
            "location": {
                "type": "string",
                "description": "Location to filter by (e.g., 'Boston', 'Nashville'). Case-insensitive partial match."
            },
            "min_capacity": {
                "type": "integer",
                "description": "Minimum venue capacity needed."
            },
            "max_capacity": {
                "type": "integer",
                "description": "Maximum venue capacity desired."
            },
            "genre": {
                "type": "string",
                "description": "Genre to filter by. Searches in the venue's genres_booked list."
            }
        }
    }
}

# Venue details tool
GET_VENUE_DETAILS_SCHEMA = {
    "name": "get_venue_details",
    "description": "Get complete profile information for a specific venue by ID. Includes description, booking contact, and typical pay range.",
    "input_schema": {
        "type": "object",
        "properties": {
            "venue_id": {
                "type": "string",
                "description": "The unique ID of the venue (e.g., 'venue_1')"
            }
        },
        "required": ["venue_id"]
    }
}

# Semantic search schemas
SEMANTIC_SEARCH_ARTISTS_SCHEMA = {
    "name": "semantic_search_artists",
    "description": "Find artists using natural language descriptions of their vibe, style, or characteristics. Use this when the user describes the type of artist they want rather than specific attributes. Examples: 'find me an artist with a chill indie vibe', 'artists with high-energy rock sound', 'intimate acoustic performers'.",
    "input_schema": {
        "type": "object",
        "properties": {
            "description": {
                "type": "string",
                "description": "Natural language description of the desired artist's vibe, style, sound, or characteristics (e.g., 'energetic rock band', 'intimate folk singer', 'electronic dance music producer'). This is the primary search parameter."
            },
            "genre": {
                "type": "string",
                "description": "Optional genre filter to narrow results (e.g., 'rock', 'jazz', 'country'). Case-insensitive."
            },
            "location": {
                "type": "string",
                "description": "Optional location/city filter (e.g., 'Boston', 'Nashville'). Case-insensitive."
            },
            "limit": {
                "type": "integer",
                "description": "Maximum number of results to return. Defaults to 10. Use lower values for quick suggestions, higher for comprehensive lists."
            }
        },
        "required": ["description"]
    }
}

SEMANTIC_SEARCH_VENUES_SCHEMA = {
    "name": "semantic_search_venues",
    "description": "Find venues using natural language descriptions of atmosphere, vibe, or characteristics. Use this when the user describes the type of venue they want rather than specific attributes. Examples: 'find me a cozy listening room', 'high-energy rock club', 'upscale jazz venue'.",
    "input_schema": {
        "type": "object",
        "properties": {
            "description": {
                "type": "string",
                "description": "Natural language description of the desired venue's vibe, atmosphere, type, or characteristics (e.g., 'intimate acoustic venue', 'large nightclub', 'legendary historic venue'). This is the primary search parameter."
            },
            "location": {
                "type": "string",
                "description": "Optional location filter (e.g., 'Boston', 'Nashville'). Case-insensitive partial match."
            },
            "min_capacity": {
                "type": "integer",
                "description": "Optional minimum venue capacity needed."
            },
            "max_capacity": {
                "type": "integer",
                "description": "Optional maximum venue capacity desired."
            },
            "genre": {
                "type": "string",
                "description": "Optional genre filter. Searches in the venue's genres_booked list."
            },
            "limit": {
                "type": "integer",
                "description": "Maximum number of results to return. Defaults to 10."
            }
        },
        "required": ["description"]
    }
}

# Tool collections for different agents
ARTIST_TOOLS = [
    SEARCH_ARTISTS_SCHEMA,
    GET_ARTIST_DETAILS_SCHEMA,
    SEMANTIC_SEARCH_ARTISTS_SCHEMA
]

VENUE_TOOLS = [
    SEARCH_VENUES_SCHEMA,
    GET_VENUE_DETAILS_SCHEMA,
    SEMANTIC_SEARCH_VENUES_SCHEMA
]

ALL_TOOLS = ARTIST_TOOLS + VENUE_TOOLS
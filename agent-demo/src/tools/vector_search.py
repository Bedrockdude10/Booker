"""Vector search query builders for MongoDB Atlas."""

from typing import Any


def build_artist_vector_search_pipeline(
    query_embedding: list[float],
    genre: str | None = None,
    location: str | None = None,
    limit: int = 10
) -> list[dict[str, Any]]:
    """Build MongoDB aggregation pipeline for artist vector search."""
    # Vector search stage (no text filters - they don't support case-insensitive matching)
    vector_search_stage = {
        "$vectorSearch": {
            "index": "artist_embedding",
            "path": "embedding",
            "queryVector": query_embedding,
            "numCandidates": limit * 20,  # Get more candidates for post-filtering
            "limit": limit * 5,  # Retrieve extra results for filtering
        }
    }

    pipeline = [vector_search_stage]

    # Add score projection first
    pipeline.append({
        "$addFields": {
            "search_score": {"$meta": "vectorSearchScore"}
        }
    })

    # Post-filter with case-insensitive matching (after vector search)
    match_filters = []
    if genre:
        # Normalize to lowercase for matching
        # Artist genres are stored lowercase, so match with lowercase query
        match_filters.append({"genres": genre.lower()})
    if location:
        # Location is stored as "City, State" - do case-insensitive substring match
        match_filters.append({"location": {"$regex": location, "$options": "i"}})

    if match_filters:
        pipeline.append({
            "$match": {"$and": match_filters} if len(match_filters) > 1 else match_filters[0]
        })

    # Limit after filtering
    pipeline.append({"$limit": limit})

    # Final projection
    pipeline.append({
        "$project": {
            "_id": 1, "id": 1, "name": 1, "genres": 1,
            "location": 1, "typical_venue_capacity": 1, "bio": 1,
            "search_score": 1
        }
    })

    return pipeline


def build_venue_vector_search_pipeline(
    query_embedding: list[float],
    location: str | None = None,
    min_capacity: int | None = None,
    max_capacity: int | None = None,
    genre: str | None = None,
    limit: int = 10
) -> list[dict[str, Any]]:
    """Build MongoDB aggregation pipeline for venue vector search."""
    # Use pre-filters for capacity (numeric comparisons work in vector search)
    vector_filters = []
    if min_capacity is not None:
        vector_filters.append({"capacity": {"$gte": min_capacity}})
    if max_capacity is not None:
        vector_filters.append({"capacity": {"$lte": max_capacity}})

    vector_search_stage = {
        "$vectorSearch": {
            "index": "venue_embedding",
            "path": "embedding",
            "queryVector": query_embedding,
            "numCandidates": limit * 20,  # Get more candidates for post-filtering
            "limit": limit * 5,  # Retrieve extra results for filtering
        }
    }

    # Add capacity filters to vector search (numeric comparisons are supported)
    if vector_filters:
        vector_search_stage["$vectorSearch"]["filter"] = (
            {"$and": vector_filters} if len(vector_filters) > 1 else vector_filters[0]
        )

    pipeline = [vector_search_stage]

    # Add score projection and lowercase genre array for matching
    pipeline.append({
        "$addFields": {
            "search_score": {"$meta": "vectorSearchScore"},
            # Create lowercase version of genres_booked for case-insensitive matching
            "genres_lower": {
                "$map": {
                    "input": "$genres_booked",
                    "as": "g",
                    "in": {"$toLower": "$$g"}
                }
            }
        }
    })

    # Post-filter for text fields
    post_filters = []
    if location:
        # Location is "City, State" format - do substring match
        post_filters.append({"location": {"$regex": location, "$options": "i"}})
    if genre:
        # Match against lowercase version of genres
        post_filters.append({"genres_lower": genre.lower()})

    if post_filters:
        pipeline.append({
            "$match": {"$and": post_filters} if len(post_filters) > 1 else post_filters[0]
        })

    # Limit after filtering
    pipeline.append({"$limit": limit})

    # Final projection (remove temp genres_lower field)
    pipeline.append({
        "$project": {
            "_id": 1, "id": 1, "name": 1, "location": 1,
            "capacity": 1, "genres_booked": 1, "venue_type": 1,
            "description": 1,
            "search_score": 1,
            "genres_lower": 0  # Exclude temp field
        }
    })

    return pipeline

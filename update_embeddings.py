"""
Update embeddings for documents that don't have them yet.

Usage:
    python update_embeddings.py
    python update_embeddings.py --collection artists
    python update_embeddings.py --collection venues
    python update_embeddings.py --dry-run  # Preview what would be updated
"""

import os
import argparse
from pymongo import MongoClient
from sentence_transformers import SentenceTransformer
from dotenv import load_dotenv

load_dotenv()


def get_embedding_model():
    """Lazy load embedding model (same as seed_database.py)."""
    return SentenceTransformer('all-mpnet-base-v2')


def create_artist_search_text(artist: dict) -> str:
    """Create search text from artist document (same logic as seed_database.py)."""
    genres = ", ".join(artist.get("genres", []))
    return f"""
    Artist: {artist.get("name", "")}
    Genres: {genres}
    Location: {artist.get("location", "")}
    Bio: {artist.get("bio", "")}
    Typical venue capacity: {artist.get("typical_venue_capacity", "")}
    Years active: {artist.get("years_active", "")}
    """.strip()


def create_venue_search_text(venue: dict) -> str:
    """Create search text from venue document (same logic as seed_database.py)."""
    genres = ", ".join(venue.get("genres_booked", []))
    return f"""
    Venue: {venue.get("name", "")}
    Location: {venue.get("location", "")}
    Capacity: {venue.get("capacity", "")}
    Venue type: {venue.get("venue_type", "")}
    Genres booked: {genres}
    Age restriction: {venue.get("ages", "")}
    Description: {venue.get("description", "")}
    Typical pay range: {venue.get("typical_pay_range", "")}
    """.strip()


def update_collection_embeddings(collection, create_text_fn, dry_run=False):
    """Update embeddings for documents missing them."""
    # Find documents without embeddings
    query = {"embedding": {"$exists": False}}
    docs_missing_embeddings = list(collection.find(query))

    if not docs_missing_embeddings:
        print(f"✅ All documents in {collection.name} already have embeddings")
        return 0

    print(f"Found {len(docs_missing_embeddings)} documents without embeddings in {collection.name}")

    if dry_run:
        print("\n🔍 DRY RUN - Would update these documents:")
        for doc in docs_missing_embeddings:
            print(f"  - {doc.get('name', doc.get('_id'))}")
        return len(docs_missing_embeddings)

    # Generate embeddings
    print("Generating embeddings...")
    model = get_embedding_model()

    for i, doc in enumerate(docs_missing_embeddings, 1):
        search_text = create_text_fn(doc)
        embedding = model.encode([search_text])[0].tolist()

        # Update document with embedding
        collection.update_one(
            {"_id": doc["_id"]},
            {"$set": {"embedding": embedding, "search_text": search_text}}
        )
        print(f"  [{i}/{len(docs_missing_embeddings)}] Updated: {doc.get('name', doc.get('_id'))}")

    print(f"✅ Updated {len(docs_missing_embeddings)} documents")
    return len(docs_missing_embeddings)


def main():
    parser = argparse.ArgumentParser(description="Update embeddings for documents missing them")
    parser.add_argument("--collection", choices=["artists", "venues", "all"], default="all",
                        help="Which collection to update (default: all)")
    parser.add_argument("--dry-run", action="store_true",
                        help="Preview changes without updating database")
    args = parser.parse_args()

    # Connect to MongoDB
    mongodb_uri = os.getenv("MONGODB_URI")
    if not mongodb_uri:
        print("❌ Error: MONGODB_URI not found in environment")
        return

    client = MongoClient(mongodb_uri)
    db = client.booker

    print(f"Connected to MongoDB (database: booker)")
    if args.dry_run:
        print("🔍 DRY RUN MODE - No changes will be made\n")

    total_updated = 0

    # Update artists
    if args.collection in ["artists", "all"]:
        print("\n--- Artists Collection ---")
        count = update_collection_embeddings(
            db.artists,
            create_artist_search_text,
            dry_run=args.dry_run
        )
        total_updated += count

    # Update venues
    if args.collection in ["venues", "all"]:
        print("\n--- Venues Collection ---")
        count = update_collection_embeddings(
            db.venues,
            create_venue_search_text,
            dry_run=args.dry_run
        )
        total_updated += count

    print(f"\n{'[DRY RUN] Would update' if args.dry_run else 'Updated'} {total_updated} total documents")

    client.close()


if __name__ == "__main__":
    main()

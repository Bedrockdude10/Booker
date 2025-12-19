"""
Seed MongoDB Atlas with mock data + embeddings for vector search.

Usage:
    pip install pymongo sentence-transformers python-dotenv
    # Make sure .env has MONGODB_URI set
    python seed_atlas.py

This script:
1. Connects to MongoDB Atlas
2. Creates 'booker' database with 'artists' and 'venues' collections
3. Generates text embeddings for semantic search (free, local model)
4. Inserts all documents with embeddings
"""

import os
from datetime import datetime
from pymongo import MongoClient
from sentence_transformers import SentenceTransformer
from dotenv import load_dotenv

# Load .env file from current directory or parent directories
load_dotenv()

# ============================================================================
# MOCK DATA (same as your agent-demo)
# ============================================================================

MOCK_ARTISTS = [
    {
        "id": "artist_1",
        "name": "The Midnight Riders",
        "genres": ["Rock", "Indie Rock", "Alternative"],
        "location": "Boston, MA",
        "bio": "High-energy rock band with a loyal following. Known for explosive live performances and original songwriting.",
        "social_links": {
            "instagram": "@midnightriders",
            "spotify": "spotify.com/artist/midnightriders"
        },
        "typical_venue_capacity": "200-500",
        "years_active": 5,
        "booking_email": "booking@midnightriders.com"
    },
    {
        "id": "artist_2",
        "name": "Sarah Chen",
        "genres": ["Folk", "Singer-Songwriter", "Acoustic"],
        "location": "Nashville, TN",
        "bio": "Intimate storyteller with haunting vocals. Perfect for listening rooms and acoustic venues.",
        "social_links": {
            "instagram": "@sarahchenmusic",
            "website": "sarahchenmusic.com"
        },
        "typical_venue_capacity": "50-200",
        "years_active": 8,
        "booking_email": "booking@sarahchenmusic.com"
    },
    {
        "id": "artist_3",
        "name": "DJ Neon Pulse",
        "genres": ["Electronic", "House", "Techno"],
        "location": "Boston, MA",
        "bio": "Genre-bending electronic producer and DJ. Brings the dance floor to life with cutting-edge beats.",
        "social_links": {
            "soundcloud": "soundcloud.com/neonpulse",
            "instagram": "@djneonpulse"
        },
        "typical_venue_capacity": "300-1000",
        "years_active": 6,
        "booking_email": "bookings@neonpulse.net"
    },
    {
        "id": "artist_4",
        "name": "The Bluegrass Collective",
        "genres": ["Bluegrass", "Country", "Americana"],
        "location": "Nashville, TN",
        "bio": "Traditional bluegrass quartet with modern sensibilities. Perfect for festivals and honky-tonks.",
        "social_links": {
            "facebook": "facebook.com/bluegrasscollective",
            "website": "bluegrasscollective.com"
        },
        "typical_venue_capacity": "100-400",
        "years_active": 12,
        "booking_email": "book@bluegrasscollective.com"
    },
    {
        "id": "artist_5",
        "name": "Velvet Underground Jazz Trio",
        "genres": ["Jazz", "Bebop", "Contemporary Jazz"],
        "location": "Boston, MA",
        "bio": "Sophisticated jazz trio with a modern edge. Great for upscale venues and jazz clubs.",
        "social_links": {
            "instagram": "@velvetjazz",
            "spotify": "spotify.com/artist/velvetjazz"
        },
        "typical_venue_capacity": "80-250",
        "years_active": 10,
        "booking_email": "contact@velvetjazz.com"
    },
    {
        "id": "artist_6",
        "name": "The Punk Revival",
        "genres": ["Punk", "Punk Rock", "Hardcore"],
        "location": "Boston, MA",
        "bio": "Fast, loud, and unapologetic. Classic punk energy with a fresh attitude.",
        "social_links": {
            "instagram": "@punkrevival",
            "bandcamp": "punkrevival.bandcamp.com"
        },
        "typical_venue_capacity": "150-400",
        "years_active": 3,
        "booking_email": "shows@punkrevival.com"
    },
    {
        "id": "artist_7",
        "name": "Luna Rodriguez",
        "genres": ["R&B", "Soul", "Neo-Soul"],
        "location": "Nashville, TN",
        "bio": "Smooth vocals meet contemporary R&B. Creates an intimate, sophisticated atmosphere.",
        "social_links": {
            "instagram": "@lunarodriguezmusic",
            "spotify": "spotify.com/artist/lunarodriguez",
            "tiktok": "@lunarodriguezmusic"
        },
        "typical_venue_capacity": "200-600",
        "years_active": 4,
        "booking_email": "mgmt@lunarodriguez.com"
    },
    {
        "id": "artist_8",
        "name": "The Heavy Hearts",
        "genres": ["Metal", "Hard Rock", "Heavy Metal"],
        "location": "Boston, MA",
        "bio": "Crushing riffs and powerful vocals. For venues that can handle high volume and energy.",
        "social_links": {
            "instagram": "@heavyheartsband",
            "youtube": "youtube.com/@heavyhearts"
        },
        "typical_venue_capacity": "300-800",
        "years_active": 7,
        "booking_email": "booking@heavyhearts.net"
    },
    {
        "id": "artist_9",
        "name": "Cosmic Country Band",
        "genres": ["Country", "Outlaw Country", "Americana"],
        "location": "Nashville, TN",
        "bio": "Traditional country with a cosmic twist. Perfect for honky-tonks and dive bars.",
        "social_links": {
            "instagram": "@cosmiccountry",
            "website": "cosmiccountryband.com"
        },
        "typical_venue_capacity": "150-500",
        "years_active": 6,
        "booking_email": "book@cosmiccountryband.com"
    },
    {
        "id": "artist_10",
        "name": "The Indie Dreamers",
        "genres": ["Indie Pop", "Dream Pop", "Indie Rock"],
        "location": "Boston, MA",
        "bio": "Dreamy melodies and shimmering guitars. Creates an ethereal live experience.",
        "social_links": {
            "instagram": "@indiedreamers",
            "spotify": "spotify.com/artist/indiedreamers"
        },
        "typical_venue_capacity": "100-300",
        "years_active": 4,
        "booking_email": "contact@indiedreamers.com"
    },
    {
        "id": "artist_11",
        "name": "The Hip-Hop Collective",
        "genres": ["Hip-Hop", "Rap", "Conscious Hip-Hop"],
        "location": "Boston, MA",
        "bio": "Thought-provoking lyrics with boom-bap beats. High-energy live shows with crowd participation.",
        "social_links": {
            "instagram": "@hiphopcollab",
            "soundcloud": "soundcloud.com/hiphopcollab"
        },
        "typical_venue_capacity": "250-700",
        "years_active": 5,
        "booking_email": "booking@hiphopcollab.com"
    },
    {
        "id": "artist_12",
        "name": "Emily & The Wildcards",
        "genres": ["Folk Rock", "Americana", "Country"],
        "location": "Nashville, TN",
        "bio": "Heartfelt storytelling with a rock edge. Full band energy with folk sensibilities.",
        "social_links": {
            "instagram": "@emilywildcards",
            "website": "emilywildcards.com"
        },
        "typical_venue_capacity": "200-500",
        "years_active": 8,
        "booking_email": "management@emilywildcards.com"
    }
]

MOCK_VENUES = [
    {
        "id": "venue_1",
        "name": "The Paradise Rock Club",
        "location": "Boston, MA",
        "capacity": 933,
        "genres_booked": ["Rock", "Indie Rock", "Alternative", "Punk"],
        "booking_contact": "booking@paradiserock.com",
        "typical_pay_range": "$800-2000",
        "venue_type": "club",
        "ages": "18+",
        "description": "Legendary Boston venue known for launching careers. Great sound system and intimate atmosphere."
    },
    {
        "id": "venue_2",
        "name": "Club Passim",
        "location": "Boston, MA",
        "capacity": 115,
        "genres_booked": ["Folk", "Singer-Songwriter", "Acoustic", "Americana"],
        "booking_contact": "booking@clubpassim.org",
        "typical_pay_range": "$300-800",
        "venue_type": "listening room",
        "ages": "all ages",
        "description": "Intimate listening room with impeccable acoustics. Perfect for singer-songwriters and acoustic acts."
    },
    {
        "id": "venue_3",
        "name": "The Bluebird Cafe",
        "location": "Nashville, TN",
        "capacity": 90,
        "genres_booked": ["Country", "Singer-Songwriter", "Americana", "Folk"],
        "booking_contact": "calendar@bluebirdcafe.com",
        "typical_pay_range": "$200-600",
        "venue_type": "listening room",
        "ages": "all ages",
        "description": "Iconic Nashville venue where songwriters shine. Known for in-the-round performances."
    },
    {
        "id": "venue_4",
        "name": "Royale Nightclub",
        "location": "Boston, MA",
        "capacity": 1200,
        "genres_booked": ["Electronic", "House", "Hip-Hop", "Pop"],
        "booking_contact": "talent@royaleboston.com",
        "typical_pay_range": "$1500-5000",
        "venue_type": "nightclub",
        "ages": "21+",
        "description": "Premier nightclub with state-of-the-art sound and lighting. Perfect for DJ sets and electronic acts."
    },
    {
        "id": "venue_5",
        "name": "The Station Inn",
        "location": "Nashville, TN",
        "capacity": 200,
        "genres_booked": ["Bluegrass", "Country", "Americana", "Folk"],
        "booking_contact": "booking@stationinn.com",
        "typical_pay_range": "$400-1000",
        "venue_type": "honky-tonk",
        "ages": "21+",
        "description": "Nashville's premier bluegrass venue. Authentic atmosphere and dedicated bluegrass fans."
    },
    {
        "id": "venue_6",
        "name": "Scullers Jazz Club",
        "location": "Boston, MA",
        "capacity": 200,
        "genres_booked": ["Jazz", "Contemporary Jazz", "Blues", "Soul"],
        "booking_contact": "info@scullersjazz.com",
        "typical_pay_range": "$500-1500",
        "venue_type": "jazz club",
        "ages": "21+",
        "description": "Upscale jazz venue with waterfront views. Dinner and drinks with world-class jazz."
    },
    {
        "id": "venue_7",
        "name": "The Sinclair",
        "location": "Boston, MA",
        "capacity": 525,
        "genres_booked": ["Indie Rock", "Alternative", "Rock", "Pop"],
        "booking_contact": "booking@sinclaircambridge.com",
        "typical_pay_range": "$1000-3000",
        "venue_type": "club",
        "ages": "18+",
        "description": "Modern venue in Harvard Square. Great production capabilities and excellent sightlines."
    },
    {
        "id": "venue_8",
        "name": "Exit/In",
        "location": "Nashville, TN",
        "capacity": 500,
        "genres_booked": ["Rock", "Indie Rock", "Alternative", "Punk", "Metal"],
        "booking_contact": "booking@exitin.com",
        "typical_pay_range": "$700-2000",
        "venue_type": "club",
        "ages": "18+",
        "description": "Historic Nashville venue with a rock and roll heart. Diverse bookings and loyal crowds."
    },
    {
        "id": "venue_9",
        "name": "The Middle East - Downstairs",
        "location": "Boston, MA",
        "capacity": 194,
        "genres_booked": ["Punk", "Metal", "Hardcore", "Alternative"],
        "booking_contact": "booking@mideastoffers.com",
        "typical_pay_range": "$400-1000",
        "venue_type": "club",
        "ages": "18+",
        "description": "Legendary underground venue. Perfect for punk, metal, and high-energy alternative acts."
    },
    {
        "id": "venue_10",
        "name": "The Basement",
        "location": "Nashville, TN",
        "capacity": 600,
        "genres_booked": ["Rock", "Indie Rock", "Alternative", "Electronic", "Hip-Hop"],
        "booking_contact": "talent@thebasementnashville.com",
        "typical_pay_range": "$800-2500",
        "venue_type": "club",
        "ages": "18+",
        "description": "Nashville's eclectic music venue. Books diverse acts and creates memorable experiences."
    },
    {
        "id": "venue_11",
        "name": "Brighton Music Hall",
        "location": "Boston, MA",
        "capacity": 380,
        "genres_booked": ["Indie Rock", "Alternative", "Pop", "Electronic"],
        "booking_contact": "booking@brightonmusichall.com",
        "typical_pay_range": "$600-1800",
        "venue_type": "club",
        "ages": "18+",
        "description": "Mid-sized venue with great energy. Perfect stepping stone for growing indie acts."
    },
    {
        "id": "venue_12",
        "name": "The Ryman Auditorium",
        "location": "Nashville, TN",
        "capacity": 2362,
        "genres_booked": ["Country", "Americana", "Folk", "Bluegrass", "Rock"],
        "booking_contact": "booking@ryman.com",
        "typical_pay_range": "$5000-15000",
        "venue_type": "theater",
        "ages": "all ages",
        "description": "The Mother Church of Country Music. Legendary acoustics and historic significance."
    },
    {
        "id": "venue_13",
        "name": "The Beehive",
        "location": "Boston, MA",
        "capacity": 150,
        "genres_booked": ["Jazz", "Soul", "R&B", "Blues"],
        "booking_contact": "events@beehiveboston.com",
        "typical_pay_range": "$400-1000",
        "venue_type": "bar",
        "ages": "21+",
        "description": "Bohemian restaurant and bar with live music. Great for jazz, soul, and R&B acts."
    },
    {
        "id": "venue_14",
        "name": "3rd and Lindsley",
        "location": "Nashville, TN",
        "capacity": 550,
        "genres_booked": ["Blues", "Rock", "Soul", "R&B", "Americana"],
        "booking_contact": "booking@3rdandlindsley.com",
        "typical_pay_range": "$800-2000",
        "venue_type": "club",
        "ages": "18+",
        "description": "Premier listening room and grill. Known for excellent sound and diverse bookings."
    }
]


# ============================================================================
# EMBEDDING GENERATION (Free, local model)
# ============================================================================

# Load model once - all-mpnet-base-v2 is the best quality, 768 dimensions
EMBEDDING_MODEL = None

def get_embedding_model():
    """Lazy load the embedding model."""
    global EMBEDDING_MODEL
    if EMBEDDING_MODEL is None:
        print("Loading embedding model (first time only, ~420MB download)...")
        EMBEDDING_MODEL = SentenceTransformer('all-mpnet-base-v2')
    return EMBEDDING_MODEL


def create_artist_search_text(artist: dict) -> str:
    """Create a rich text representation for embedding."""
    genres = ", ".join(artist["genres"])
    return f"""
    Artist: {artist["name"]}
    Genres: {genres}
    Location: {artist["location"]}
    Bio: {artist["bio"]}
    Typical venue capacity: {artist["typical_venue_capacity"]}
    Years active: {artist["years_active"]}
    """.strip()


def create_venue_search_text(venue: dict) -> str:
    """Create a rich text representation for embedding."""
    genres = ", ".join(venue["genres_booked"])
    return f"""
    Venue: {venue["name"]}
    Location: {venue["location"]}
    Capacity: {venue["capacity"]}
    Venue type: {venue["venue_type"]}
    Genres booked: {genres}
    Age restriction: {venue["ages"]}
    Description: {venue["description"]}
    Typical pay range: {venue["typical_pay_range"]}
    """.strip()


def generate_embeddings(texts: list[str]) -> list[list[float]]:
    """Generate embeddings using local Sentence Transformers model."""
    model = get_embedding_model()
    embeddings = model.encode(texts, convert_to_numpy=True)
    return [emb.tolist() for emb in embeddings]


# ============================================================================
# MAIN SEEDING LOGIC
# ============================================================================

def seed_database():
    """Main function to seed MongoDB Atlas with mock data + embeddings."""
    
    # Get connection string from environment
    mongodb_uri = os.environ.get("MONGODB_URI")
    if not mongodb_uri:
        raise ValueError("MONGODB_URI environment variable not set")
    
    # Initialize MongoDB client
    print("Connecting to MongoDB Atlas...")
    mongo_client = MongoClient(mongodb_uri)
    
    # Select database
    db = mongo_client["booker"]
    
    # Drop existing collections (fresh start)
    print("Dropping existing collections...")
    db.artists.drop()
    db.venues.drop()
    
    # -------------------------------------------------------------------------
    # SEED ARTISTS
    # -------------------------------------------------------------------------
    print(f"Processing {len(MOCK_ARTISTS)} artists...")
    
    # Generate search texts
    artist_texts = [create_artist_search_text(a) for a in MOCK_ARTISTS]
    
    # Generate embeddings in batch (free, local)
    print("Generating artist embeddings...")
    artist_embeddings = generate_embeddings(artist_texts)
    
    # Prepare documents with embeddings
    artist_docs = []
    for artist, search_text, embedding in zip(MOCK_ARTISTS, artist_texts, artist_embeddings):
        doc = {
            **artist,
            "search_text": search_text,
            "embedding": embedding,
            "created_at": datetime.utcnow(),
            "updated_at": datetime.utcnow()
        }
        artist_docs.append(doc)
    
    # Insert
    result = db.artists.insert_many(artist_docs)
    print(f"Inserted {len(result.inserted_ids)} artists")
    
    # -------------------------------------------------------------------------
    # SEED VENUES
    # -------------------------------------------------------------------------
    print(f"Processing {len(MOCK_VENUES)} venues...")
    
    # Generate search texts
    venue_texts = [create_venue_search_text(v) for v in MOCK_VENUES]
    
    # Generate embeddings in batch (free, local)
    print("Generating venue embeddings...")
    venue_embeddings = generate_embeddings(venue_texts)
    
    # Prepare documents with embeddings
    venue_docs = []
    for venue, search_text, embedding in zip(MOCK_VENUES, venue_texts, venue_embeddings):
        doc = {
            **venue,
            "search_text": search_text,
            "embedding": embedding,
            "created_at": datetime.utcnow(),
            "updated_at": datetime.utcnow()
        }
        venue_docs.append(doc)
    
    # Insert
    result = db.venues.insert_many(venue_docs)
    print(f"Inserted {len(result.inserted_ids)} venues")
    
    # -------------------------------------------------------------------------
    # CREATE INDEXES
    # -------------------------------------------------------------------------
    print("Creating standard indexes...")
    
    # Artists indexes
    db.artists.create_index("id", unique=True)
    db.artists.create_index("location")
    db.artists.create_index("genres")
    
    # Venues indexes
    db.venues.create_index("id", unique=True)
    db.venues.create_index("location")
    db.venues.create_index("genres_booked")
    db.venues.create_index("capacity")
    
    print("Standard indexes created")
    
    # -------------------------------------------------------------------------
    # VECTOR SEARCH INDEX INSTRUCTIONS
    # -------------------------------------------------------------------------
    print("\n" + "="*70)
    print("VECTOR SEARCH INDEX SETUP")
    print("="*70)
    print("""
Vector search indexes must be created in Atlas UI or via Atlas CLI.

Go to: Atlas → Your Cluster → Search → Create Search Index → JSON Editor

For ARTISTS collection, create index named 'artist_embedding':
{
  "fields": [
    {
      "type": "vector",
      "path": "embedding",
      "numDimensions": 768,
      "similarity": "cosine"
    },
    {
      "type": "filter",
      "path": "location"
    },
    {
      "type": "filter",
      "path": "genres"
    }
  ]
}

For VENUES collection, create index named 'venue_embedding':
{
  "fields": [
    {
      "type": "vector",
      "path": "embedding",
      "numDimensions": 768,
      "similarity": "cosine"
    },
    {
      "type": "filter",
      "path": "location"
    },
    {
      "type": "filter",
      "path": "genres_booked"
    },
    {
      "type": "filter",
      "path": "capacity"
    }
  ]
}
    """)
    print("="*70)
    
    print("\n✅ Seeding complete!")
    print(f"   - {len(MOCK_ARTISTS)} artists with embeddings")
    print(f"   - {len(MOCK_VENUES)} venues with embeddings")
    print("\nNext steps:")
    print("1. Create vector search indexes in Atlas UI (see above)")
    print("2. Add semantic search endpoint to your Go backend")
    print("3. Connect agent-demo to the API")


if __name__ == "__main__":
    seed_database()
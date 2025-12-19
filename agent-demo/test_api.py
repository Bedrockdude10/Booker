#!/usr/bin/env python3
"""Test the API connection. Run: python test_api.py"""

from src.tools.api_client import BookerClient

def main():
    client = BookerClient()
    print(f"Testing: {client.base_url}\n")
    
    # Health check
    print("1. Health check:", "✅" if client.health() else "❌")
    
    # Search all artists
    artists = client.search_artists()
    print(f"2. All artists: {len(artists)} found")
    if artists:
        print(f"   First: {artists[0].get('name', 'N/A')}")
    
    # Search with filter
    rock = client.search_artists(genres="rock")
    print(f"3. Rock artists: {len(rock) if rock else 0} found")
    
    # Get by ID (if we have artists)
    if artists:
        aid = artists[0].get("_id", {})
        aid = aid.get("$oid", str(aid)) if isinstance(aid, dict) else str(aid)
        detail = client.get_artist(aid)
        print(f"4. Get by ID: {'✅' if detail else '❌'}")
    
    client.close()
    print("\n✅ Connection working!")

if __name__ == "__main__":
    main()
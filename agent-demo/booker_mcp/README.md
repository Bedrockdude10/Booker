# Booker MCP Server

MCP (Model Context Protocol) server that exposes Booker's artist-venue matching tools to Claude Desktop and other MCP-compatible clients.

## What This Server Does

This server makes Booker's four main tools available via the MCP protocol:

1. **search_artists** - Search for artists by genre, location, and venue capacity preferences
2. **search_venues** - Search venues by location, capacity range, and genres booked
3. **get_artist_details** - Get complete profile information for a specific artist
4. **get_venue_details** - Get complete profile information for a specific venue

## Running the Server Standalone

To test the server standalone:

```bash
cd agent-demo
python -m booker_mcp.server
```

The server will start and listen for MCP requests via stdio.

## Configuring Claude Desktop

To add this server to Claude Desktop, edit your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add the following to the `mcpServers` section:

```json
{
  "mcpServers": {
    "booker": {
      "command": "python",
      "args": [
        "-m",
        "booker_mcp.server"
      ],
      "cwd": "/absolute/path/to/agent-demo"
    }
  }
}
```

Replace `/absolute/path/to/agent-demo` with the actual absolute path to your agent-demo directory.

See [claude_desktop_config.json](./claude_desktop_config.json) for a complete example.

## Available Tools

### search_artists
Search for artists by genre, location, and capacity preferences.

**Parameters:**
- `genre` (optional): Genre to filter by (e.g., 'Rock', 'Jazz', 'Country')
- `location` (optional): Location to filter by (e.g., 'Boston', 'Nashville')
- `max_venue_capacity` (optional): Maximum venue capacity the artist typically plays

**Returns:** List of matching artists with basic information

### search_venues
Search venues by location, capacity range, and genres booked.

**Parameters:**
- `location` (optional): Location to filter by
- `min_capacity` (optional): Minimum venue capacity needed
- `max_capacity` (optional): Maximum venue capacity desired
- `genre` (optional): Genre to filter by

**Returns:** List of matching venues with basic information

### get_artist_details
Get complete profile information for a specific artist by ID.

**Parameters:**
- `artist_id` (required): The unique ID of the artist (e.g., 'artist_1')

**Returns:** Complete artist profile including bio, contact info, and social links

### get_venue_details
Get complete profile information for a specific venue by ID.

**Parameters:**
- `venue_id` (required): The unique ID of the venue (e.g., 'venue_1')

**Returns:** Complete venue profile including description, booking contact, and pay range

## Implementation Notes

- The server uses the existing tool implementations from `src/tools/registry.py`
- Tool schemas are imported from `src/tools/schemas.py` and converted to MCP format
- All results are returned as JSON-formatted text
- The server uses async/await and runs via stdio transport

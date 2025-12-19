"""MCP server for Booker artist-venue matching tools."""

import asyncio
import json
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent

# Import existing tool infrastructure
from src.tools.registry import execute_tool
from src.tools.schemas import (
    SEARCH_ARTISTS_SCHEMA,
    SEARCH_VENUES_SCHEMA,
    GET_ARTIST_DETAILS_SCHEMA,
    GET_VENUE_DETAILS_SCHEMA
)


def convert_schema_to_mcp(tool_schema: dict) -> Tool:
    """Convert our tool schema format to MCP Tool format.

    Our schemas use "input_schema" but MCP uses "inputSchema".
    """
    return Tool(
        name=tool_schema["name"],
        description=tool_schema["description"],
        inputSchema=tool_schema["input_schema"]
    )


# Create the MCP server
server = Server("booker")


@server.list_tools()
async def list_tools() -> list[Tool]:
    """Return list of available Booker tools."""
    return [
        convert_schema_to_mcp(SEARCH_ARTISTS_SCHEMA),
        convert_schema_to_mcp(SEARCH_VENUES_SCHEMA),
        convert_schema_to_mcp(GET_ARTIST_DETAILS_SCHEMA),
        convert_schema_to_mcp(GET_VENUE_DETAILS_SCHEMA)
    ]


@server.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    """Execute a Booker tool and return the results.

    This routes to the existing execute_tool() function in registry.py.
    """
    # Execute the tool using existing infrastructure
    result = execute_tool(name, arguments)

    # Return result as JSON-formatted text
    return [TextContent(
        type="text",
        text=json.dumps(result, indent=2)
    )]


async def main():
    """Run the MCP server using stdio transport."""
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            server.create_initialization_options()
        )


if __name__ == "__main__":
    asyncio.run(main())

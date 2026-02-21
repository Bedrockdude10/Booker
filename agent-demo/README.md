# Booker: Artist-Venue Matching Multi-Agent System

A production-grade multi-agent system built with Claude that matches artists with venues using intelligent orchestration, memory management, semantic search, and comprehensive observability.

## Live Demo

🎵 **[Try the live demo](https://booker-hjbdoiu5fobsifeuwwnpsg.streamlit.app)**

## Quick Start

### Prerequisites
- Python 3.10+
- Anthropic API key ([Get one here](https://console.anthropic.com/settings/keys))
- MongoDB Atlas account (for vector search features)

### Installation

```bash
# 1. Navigate to the project
cd agent-demo

# 2. Install dependencies
pip install -e .
# Or use requirements from root:
pip install -r ../requirements.txt

# 3. Set up environment variables
cp .env.example .env
# Edit .env and add:
#   - CLAUDE_API_KEY
#   - MONGODB_URI (for semantic search)
#   - BOOKER_API_URL (optional, defaults to production)

# 4. Test the system (optional but recommended)
python test_system.py

# 5. Launch the application
streamlit run app/main.py
```

The application will open at `http://localhost:8501`

### First Steps

1. **Try a simple query**: "Find me rock venues in Boston"
2. **Try semantic search**: "Find artists with a dreamy indie vibe"
3. **Toggle observability**: Click "Show Traces & Metrics" in the sidebar

## Architecture

### System Overview

```
User Query → Coordinator Agent → Specialist Agents → Tools → Response
                 ↓                      ↓               ↓
             Routing              Tool Execution    Data Access
                 ↓                      ↓               ↓
           Governance           Memory Systems      Metrics
                                                        ↓
                                              ┌─────────┴─────────┐
                                              │                   │
                                         Go Backend        MongoDB Atlas
                                         (Artists)       (Vector Search)
```

### Multi-Agent Architecture

| Agent | Role | Tools |
|-------|------|-------|
| **Coordinator** | Routes requests to specialists | `route_to_agent` |
| **Artist Discovery** | Finds and ranks artists | `search_artists`, `get_artist_details`, `semantic_search_artists` |
| **Venue Matching** | Finds and scores venues | `search_venues`, `get_venue_details`, `semantic_search_venues` |
| **Booking Advisor** | Synthesizes recommendations | All tools |

### Data Sources

| Entity | Source | Features |
|--------|--------|----------|
| **Artists** | Go Backend + MongoDB Atlas | REST API filtering by name/genre/location, vector search by vibe/description |
| **Venues** | Placeholder data + MongoDB Atlas | Mock data for filter search, vector search for semantic queries |

### Directory Structure

```
agent-demo/
├── src/
│   ├── agents/              # Agent implementations
│   │   ├── base.py          # Abstract base agent
│   │   ├── coordinator.py   # Request routing
│   │   ├── artist_discovery.py
│   │   ├── venue_matching.py
│   │   └── booking_advisor.py
│   ├── memory/              # Memory systems
│   │   ├── conversation.py  # Chat history
│   │   ├── working.py       # Intermediate results
│   │   └── preferences.py   # User preferences
│   ├── orchestration/       # Agent coordination
│   │   └── executor.py
│   ├── observability/       # Tracing & monitoring
│   │   ├── tracer.py
│   │   ├── logger.py
│   │   └── metrics.py
│   ├── governance/          # Safety & compliance
│   │   ├── coordinator.py
│   │   ├── nemo_wrapper.py
│   │   ├── pii_protection.py
│   │   ├── cost_control.py
│   │   └── audit.py
│   ├── tools/               # Tool implementations
│   │   ├── api_client.py    # Go backend HTTP client
│   │   ├── mongo_client.py  # MongoDB Atlas client
│   │   ├── embeddings.py    # Sentence transformers
│   │   ├── vector_search.py # Atlas vector search pipelines
│   │   ├── registry.py      # Tool execution
│   │   └── schemas.py       # Tool definitions
│   ├── models/              # Data models
│   │   ├── message.py
│   │   └── trace.py
│   └── config/              # Configuration
│       ├── settings.py
│       └── prompts.py
├── app/                     # Streamlit UI
│   ├── main.py
│   └── components/
│       ├── trace_viewer.py
│       └── metrics_panel.py
├── booker_mcp/              # MCP Server
│   ├── server.py
│   ├── README.md
│   └── claude_desktop_config.json
├── config/                  # Configuration files
│   ├── governance.yaml
│   └── governance.dev.yaml
└── tests/                   # Test suite
```

## Features

### Semantic Search (Vector Embeddings)

Uses MongoDB Atlas Vector Search with `all-mpnet-base-v2` embeddings (768 dimensions):

```python
# Find artists by vibe/description
semantic_search_artists(
    description="dreamy indie rock with folk influences",
    location="Boston",
    limit=5
)

# Find venues by atmosphere
semantic_search_venues(
    description="intimate listening room with great acoustics",
    min_capacity=100,
    max_capacity=300
)
```

The vector search pipelines use post-filtering for case-insensitive text matching and pre-filtering for numeric fields (capacity ranges), with oversampling (`numCandidates = limit * 20`) to ensure quality results after filtering.

### Memory Management
- **Conversation Memory**: Maintains chat history across sessions (50 messages default)
- **Working Memory**: Tracks intermediate results during multi-step tasks
- **Preference Memory**: Learns and persists user preferences (genres, locations, capacity)

### Observability
- **Execution Tracing**: Detailed traces of all agent interactions via OpenTelemetry-compatible tracer
- **Structured Logging**: JSON-formatted logs with context
- **Performance Metrics**: Token usage and duration tracking per agent
- **Real-time Visualization**: Live trace and metrics in split-view UI

### MCP Server Integration

The `booker_mcp/` directory provides an MCP (Model Context Protocol) server for Claude Desktop integration. Currently exposes the four core tools (`search_artists`, `search_venues`, `get_artist_details`, `get_venue_details`):

```json
{
  "mcpServers": {
    "booker": {
      "command": "python",
      "args": ["-m", "booker_mcp.server"],
      "cwd": "/path/to/agent-demo"
    }
  }
}
```

See [booker_mcp/README.md](booker_mcp/README.md) for setup details.

### Production Governance (Optional)

- **Content Safety**: NeMo Guardrails for jailbreak detection and topic control
- **Cost Controls**: Budget limits and rate limiting with pyrate-limiter (global, per-session, per-user)
- **Audit Logging**: Complete compliance trail for all requests (`logs/audit/`)
- **Tool Validation**: Ensures agents only use authorized tools

## Configuration

### Environment Variables

```bash
# Required
CLAUDE_API_KEY=your_key_here

# Backend Connection
BOOKER_API_URL=https://booker-65350421664.europe-west1.run.app  # Default
MONGODB_URI=mongodb+srv://...  # For semantic search

# Model Settings
DEFAULT_MODEL=claude-sonnet-4-5-20250929
MAX_TOKENS=4096

# Memory
CONVERSATION_HISTORY_LIMIT=50

# Observability
ENABLE_TRACING=true
LOG_LEVEL=INFO

# Agent Behavior
MAX_AGENT_ITERATIONS=10

# Governance (optional)
ENABLE_GOVERNANCE=false
GOVERNANCE_CONFIG=config/governance.yaml
```

## Tools Reference

### Artist Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `search_artists` | Filter-based search via Go backend | `name`, `genre`, `location`, `max_venue_capacity` |
| `get_artist_details` | Get full profile via Go backend | `artist_id` (required) |
| `semantic_search_artists` | Vibe/description search via MongoDB Atlas | `description` (required), `genre`, `location`, `limit` |

### Venue Tools

| Tool | Description | Parameters |
|------|-------------|------------|
| `search_venues` | Filter-based search (placeholder data) | `location`, `min_capacity`, `max_capacity`, `genre` |
| `get_venue_details` | Get full profile (placeholder data) | `venue_id` (required) |
| `semantic_search_venues` | Atmosphere/vibe search via MongoDB Atlas | `description` (required), `location`, `min_capacity`, `max_capacity`, `genre`, `limit` |

## Development

### Adding a New Agent

```python
from src.agents.base import BaseAgent, AgentResponse

class MyNewAgent(BaseAgent):
    def __init__(self, client=None):
        super().__init__(
            name="my_new_agent",
            system_prompt="Your prompt here",
            tools=[],  # Add tool schemas
            client=client
        )

    def process(self, user_message, context=None, conversation_history=None):
        # Your logic here
        return AgentResponse(content="Response", tokens_in=0, tokens_out=0)
```

### Adding New Tools

1. Define schema in `src/tools/schemas.py`
2. Implement function in `src/tools/registry.py`
3. Register in `_TOOL_REGISTRY`

### Running Tests

```bash
pytest tests/
```

## Technology Stack

- **LLM**: Claude Sonnet 4.5 via Anthropic API
- **Backend**: Go with Chi router, deployed on Google Cloud Run
- **Database**: MongoDB Atlas with Vector Search
- **Embeddings**: Sentence Transformers (`all-mpnet-base-v2`, 768 dimensions)
- **UI**: Streamlit (deployed on Streamlit Cloud)
- **Protocol**: MCP (Model Context Protocol) for Claude Desktop
- **Observability**: OpenTelemetry-compatible tracing
- **Governance** (optional): NeMo Guardrails, pyrate-limiter

## Troubleshooting

### API errors?
- Check `.env` has valid `CLAUDE_API_KEY`
- Verify backend is reachable: `curl https://booker-65350421664.europe-west1.run.app/health`

### Semantic search not working?
- Ensure `MONGODB_URI` is set
- Check vector search index exists in Atlas (`artist_embedding`, `venue_embedding`)
- The embedding model (`all-mpnet-base-v2`) loads lazily on first query — expect a brief delay

### Import errors?
- Run from `agent-demo` directory
- Install with `pip install -e .`

### Governance issues?
```bash
export ENABLE_GOVERNANCE=false  # Disable entirely
```

## Known Issues

- **Agent prompts don't reference semantic search tools**: The system prompts in `prompts.py` only list the core filter-based tools. Agents receive semantic search tools in their tool list but lack explicit prompt guidance on when to choose semantic vs. filter search. This can cause inconsistent tool selection for subjective/vibe-based queries.
- **MCP server exposes 4 of 6 tools**: The MCP server (`booker_mcp/server.py`) currently only exposes the four core tools, not the two semantic search tools.
- **Venue data is placeholder**: `search_venues` and `get_venue_details` use hardcoded mock data. Semantic venue search hits MongoDB Atlas, which may have different data than the placeholder list.

## Roadmap

### ✅ Completed
- Multi-agent orchestration with coordinator pattern
- Go backend with MongoDB Atlas integration
- Semantic search with vector embeddings (artists and venues)
- MCP server for Claude Desktop
- Memory systems (conversation, working, preferences)
- Observability infrastructure
- Streamlit Cloud deployment
- Bandcamp artist scraping pipeline

### 🔄 In Progress
- Enhanced agent prompts for semantic search tool selection
- Venue endpoints in Go backend

### 📋 Planned
- Expose semantic search tools via MCP server
- Real venue data integration
- User authentication
- Preference learning improvements
- Additional specialist agents

## License

MIT License - see LICENSE file for details

---

Built with Claude Sonnet 4.5 | [Go Backend](../backend) | [MongoDB Atlas](https://cloud.mongodb.com)

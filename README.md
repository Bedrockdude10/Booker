# Booker

**Booker** is a multi-agent AI system that matches artists with venues. It demonstrates production patterns for building agentic applications with Claude: coordinator-based routing, specialist agents, semantic search, memory management, an MCP server, and an observability layer — all wired to a real Go REST backend backed by MongoDB Atlas.

## Live Demo

[Try the Streamlit app](https://booker-hjbdoiu5fobsifeuwwnpsg.streamlit.app)

---

## Architecture

```
User Query
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│                   Coordinator Agent                      │
│   Routes requests to the right specialist based on       │
│   intent (artist discovery, venue matching, or advice)   │
└──────────┬──────────────────┬───────────────────────────┘
           │                  │
           ▼                  ▼
  ┌─────────────────┐  ┌─────────────────────┐
  │ Artist Discovery │  │   Venue Matching     │
  │ Agent            │  │   Agent              │
  └────────┬─────────┘  └──────────┬──────────┘
           │                       │
           └──────────┬────────────┘
                      ▼
           ┌────────────────────┐
           │  Booking Advisor   │
           │  Agent             │
           └────────┬───────────┘
                    │
         ┌──────────┴──────────┐
         ▼                     ▼
  ┌─────────────┐     ┌──────────────────┐
  │  Go Backend  │     │  MongoDB Atlas   │
  │  REST API    │     │  Vector Search   │
  │  (artists)   │     │  (embeddings)    │
  └─────────────┘     └──────────────────┘
```

### Components

| Component | Location | Description |
|-----------|----------|-------------|
| **Multi-agent system** | `agent-demo/` | Python — coordinator + 3 specialist agents |
| **MCP server** | `agent-demo/booker_mcp/` | Exposes Booker tools to Claude Desktop |
| **Streamlit UI** | `agent-demo/app/` | Chat interface with live trace viewer |
| **Go backend** | `backend/` | REST API with MongoDB Atlas integration |
| **Mobile app** | `frontend-mobile/` | React Native (Expo) |

### Agent Roles

| Agent | Responsibility | Key Tools |
|-------|----------------|-----------|
| **Coordinator** | Classifies intent, routes to specialists | `route_to_agent` |
| **Artist Discovery** | Finds and ranks artists | `search_artists`, `semantic_search_artists`, `get_artist_details` |
| **Venue Matching** | Finds and scores venues | `search_venues`, `semantic_search_venues`, `get_venue_details` |
| **Booking Advisor** | Synthesizes cross-agent recommendations | All tools |

---

## Key Patterns Demonstrated

- **Coordinator routing** — a lightweight coordinator agent classifies intent and delegates to specialist agents, keeping each agent focused and testable
- **MCP server** — `booker_mcp/` exposes the same tools to Claude Desktop via the Model Context Protocol
- **Semantic search** — MongoDB Atlas Vector Search with 768-dimension `all-mpnet-base-v2` embeddings for vibe/description queries ("find artists with a dreamy indie sound")
- **Three-tier memory** — conversation history, working memory for intermediate results, and persistent preference learning
- **Observability** — OpenTelemetry-compatible tracing with real-time trace viewer in the UI
- **Governance layer** (optional) — NeMo Guardrails, token budget controls, and audit logging; disabled by default

---

## Quickstart

### Prerequisites

- Python 3.10+
- Anthropic API key ([console.anthropic.com](https://console.anthropic.com/settings/keys))
- MongoDB Atlas cluster with vector search enabled ([cloud.mongodb.com](https://cloud.mongodb.com))
- Go 1.21+ (only if running the backend locally)

### 1. Run the agent demo (Streamlit UI)

```bash
cd agent-demo

pip install -r ../requirements.txt

cp .env.example .env
# Fill in CLAUDE_API_KEY and MONGODB_URI

streamlit run app/main.py
# Opens at http://localhost:8501
```

### 2. Run the Go backend

```bash
cd backend

cp ../.env.example .env
# Fill in MONGODB_URI

go run .
# Starts at http://localhost:8080
```

### 3. Seed the database

```bash
# From repo root — populates MongoDB with artist and venue data
pip install -r requirements.txt
cp .env.example .env  # fill in MONGODB_URI
python seed_database.py
python update_embeddings.py
```

### 4. Connect to Claude Desktop via MCP

See [agent-demo/booker_mcp/README.md](agent-demo/booker_mcp/README.md) for setup.

---

## Project Structure

```
Booker/
├── agent-demo/              # Multi-agent system (Python)
│   ├── src/
│   │   ├── agents/          # Coordinator + 3 specialist agents
│   │   ├── orchestration/   # Executor (runs the agent loop)
│   │   ├── tools/           # Tool schemas, registry, API + MongoDB clients
│   │   ├── memory/          # Conversation, working, and preference memory
│   │   ├── observability/   # Tracer, logger, metrics
│   │   ├── governance/      # NeMo, cost control, audit (optional)
│   │   ├── models/          # Shared data models
│   │   └── config/          # Settings and system prompts
│   ├── app/                 # Streamlit UI + components
│   ├── booker_mcp/          # MCP server for Claude Desktop
│   ├── config/              # governance.yaml
│   └── tests/
├── backend/                 # Go REST API (Chi router + MongoDB Atlas)
│   ├── handlers/            # artists, discovery, recommendations, auth
│   ├── domain/              # Domain models
│   └── docs/                # Swagger/OpenAPI specs
├── frontend-mobile/         # React Native app (Expo)
├── seed_database.py         # Populate MongoDB with sample data
├── update_embeddings.py     # Generate/refresh vector embeddings
└── requirements.txt         # Python dependencies
```

---

## Configuration

Copy the appropriate `.env.example` and fill in your values:

| Variable | Required | Description |
|----------|----------|-------------|
| `CLAUDE_API_KEY` | Yes | Anthropic API key |
| `MONGODB_URI` | Yes | MongoDB Atlas connection string |
| `BOOKER_API_URL` | No | Go backend URL (default: `http://localhost:8080`) |
| `ENABLE_GOVERNANCE` | No | Enable NeMo + cost controls (default: `false`) |
| `ENABLE_TRACING` | No | Enable OpenTelemetry tracing (default: `true`) |

See `agent-demo/.env.example` for the full list.

---

## Technology Stack

| Layer | Technology |
|-------|-----------|
| LLM | Claude Sonnet 4.5 (Anthropic API) |
| Agent framework | Claude Agent SDK (Python) |
| Backend | Go, Chi router, Google Cloud Run |
| Database | MongoDB Atlas + Vector Search |
| Embeddings | Sentence Transformers (`all-mpnet-base-v2`) |
| UI | Streamlit |
| Protocol | MCP (Model Context Protocol) |
| Observability | OpenTelemetry |
| Governance | NeMo Guardrails, pyrate-limiter |

---

## License

MIT License — see [LICENSE](LICENSE) for details.

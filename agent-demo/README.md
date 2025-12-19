# Artist-Venue Matching Multi-Agent System

A production-grade multi-agent system built with Claude that matches artists with venues using intelligent orchestration, memory management, and comprehensive observability.

## Quick Start

### Prerequisites
- Python 3.10+
- Anthropic API key ([Get one here](https://console.anthropic.com/settings/keys))

### Installation

```bash
# 1. Navigate to the project
cd agent-demo

# 2. Install dependencies
pip install -r requirements.txt

# 3. Download spaCy model for governance features
python -m spacy download en_core_web_lg

# 4. Set up environment variables
cp .env.example .env
# Edit .env and add your CLAUDE_API_KEY

# 5. Test the system (optional but recommended)
python test_system.py

# 6. Launch the application
streamlit run app/main.py
```

The application will open at `http://localhost:8501`

### First Steps

1. **Try a simple query**: "Find me rock venues in Boston"
2. **Toggle observability**: Click "Show Traces & Metrics" in the sidebar
3. **Explore different queries**:
   - "Show me folk artists in Nashville"
   - "I need an indie rock venue with 200-300 capacity"
   - "Recommend some artist-venue pairings for jazz"

## Complex Test Queries

These queries thoroughly exercise the multi-agent system's capabilities:

### Query 1: Multi-City Tour Planning (High Complexity)
```
I manage a folk rock band from Nashville that typically draws 200-500 people. We're planning our first Boston tour next month and need 2-3 venue recommendations. What venues would be good matches and what pay range should we expect? Also, are there any local Boston folk artists we should know about for potential collaboration?
```
**Tests**: Multi-agent routing, cross-city synthesis, multiple tool calls, practical booking advice

### Query 2: Festival Planning (Very High Complexity)
```
I'm planning a jazz festival in Boston and need help pairing artists with venues. I want to book 3 different sized venues - a small intimate listening room (under 150 capacity), a mid-sized jazz club (200-400), and a larger theater (500+). Can you recommend appropriate jazz artists for each venue based on their typical draw, and explain why each pairing makes sense?
```
**Tests**: Complex filtering across capacity tiers, multiple artist-venue pairings, Booking Advisor synthesis

### Query 3: Venue Manager Perspective (High Complexity)
```
I'm the booking manager at Paradise Rock Club in Boston. We have an opening next Friday and typically book rock and indie acts. Our capacity is around 900. I want to book a local band that can draw well and fits our vibe. Who would you recommend and what should I expect to pay them? Also, which other Boston venues should I watch as competitors?
```
**Tests**: Reverse matching (venue → artists), local filtering, pay range estimation, competitive analysis

### Query 4: Conversation Memory Test (Multi-Turn)
**Turn 1:**
```
Show me electronic music venues in Boston
```
**Turn 2:**
```
Which of those would work for a DJ who typically draws 500-1000 people?
```
**Turn 3:**
```
Perfect! Can you give me the booking contact and typical pay range for the best match?
```
**Tests**: Conversation memory across turns, context retention, progressive refinement

### Query 5: Edge Case Handling (Complexity + Error Handling)
```
I need a metal venue in Nashville with capacity under 200. If that doesn't exist in your database, what would be the closest alternative? Could I book a metal band at a different type of venue, and would that work?
```
**Tests**: Handling queries with no exact matches, creative problem-solving, fallback recommendations

### Query 6: Budget-Constrained Booking (High Complexity)
```
I'm organizing a small country music showcase in Nashville with a tight budget. I need 2-3 venues that book country/americana acts with capacities under 300, and local Nashville artists who would be appropriate. What's the most affordable combination you can recommend while still maintaining quality?
```
**Tests**: Budget awareness, multiple venue/artist recommendations, location/genre filtering, value optimization

### Query 7: Genre Crossover Analysis (Medium-High Complexity)
```
I have a band that plays indie rock with folk influences. They're from Boston and typically play to 150-300 people. What venues in Boston would work for this kind of crossover sound, and would you recommend positioning them as indie rock or folk for booking purposes?
```
**Tests**: Multi-genre reasoning, genre strategy advice, capacity matching, booking positioning

## Features

### Multi-Agent Architecture
- **Coordinator Agent**: Intelligently routes requests to specialist agents
- **Artist Discovery Agent**: Searches and ranks artists by criteria
- **Venue Matching Agent**: Finds and scores venues for fit
- **Booking Advisor Agent**: Synthesizes information and provides recommendations

### Memory Management
- **Conversation Memory**: Maintains chat history across sessions (50 messages default)
- **Working Memory**: Tracks intermediate results during multi-step tasks
- **Preference Memory**: Learns and persists user preferences

### Observability
- **Execution Tracing**: Detailed traces of all agent interactions
- **Structured Logging**: JSON-formatted logs with context
- **Performance Metrics**: Token usage and duration tracking per agent
- **Real-time Visualization**: Live trace and metrics in split-view UI

### Production-Ready Governance
- **Content Safety**: NeMo Guardrails for jailbreak detection and topic control
<!-- - **PII Protection**: Microsoft Presidio for detecting and anonymizing sensitive data -->
- **Cost Controls**: Budget limits and rate limiting with pyrate-limiter
- **Audit Logging**: Complete compliance trail for all requests
- **Tool Validation**: Ensures agents only use authorized tools

## Architecture

### System Overview

```
User Query → Coordinator Agent → Specialist Agents → Tools → Response
                 ↓                      ↓               ↓
             Routing              Tool Execution    Data Access
                 ↓                      ↓               ↓
           Governance           Memory Systems      Metrics
```

### Directory Structure

```
agent-demo/
├── src/
│   ├── agents/          # Agent implementations
│   │   ├── base.py      # Abstract base agent
│   │   ├── coordinator.py
│   │   ├── artist_discovery.py
│   │   ├── venue_matching.py
│   │   └── booking_advisor.py
│   ├── memory/          # Memory systems
│   │   ├── conversation.py
│   │   ├── working.py
│   │   └── preferences.py
│   ├── orchestration/   # Agent coordination
│   │   └── executor.py
│   ├── observability/   # Tracing & monitoring
│   │   ├── tracer.py
│   │   ├── logger.py
│   │   └── metrics.py
│   ├── governance/      # Safety & compliance
│   │   ├── coordinator.py
│   │   ├── nemo_wrapper.py
│   │   ├── pii_protection.py
│   │   ├── cost_control.py
│   │   └── audit.py
│   ├── tools/           # Tool definitions
│   │   ├── registry.py
│   │   └── schemas.py
│   ├── models/          # Data models
│   │   ├── message.py
│   │   └── trace.py
│   └── config/          # Configuration
│       ├── settings.py
│       └── prompts.py
├── app/                 # Streamlit UI
│   ├── main.py
│   └── components/
│       ├── trace_viewer.py
│       └── metrics_panel.py
├── data/                # Mock data
│   └── mock_data.py
├── config/              # Configuration files
│   ├── governance.yaml
│   └── governance.dev.yaml
└── tests/               # Test suite
```

### Request Flow

```
User Input
    ↓
[Layer 1] Input Validation & Content Moderation (NeMo)
    ↓
[Layer 2] Budget Control & Rate Limiting (pyrate-limiter)
    ↓
[Layer 3] PII Detection & Anonymization (Presidio)
    ↓
[Layer 4] Agent Execution with Monitoring
    │
    ├─▶ Add to conversation memory
    ├─▶ Build context (preferences, etc.)
    ├─▶ Coordinator analyzes intent and routes
    ├─▶ Specialist agent executes with tools
    └─▶ Record metrics & traces
    ↓
[Layer 5] Output Filtering & PII Protection
    ↓
[Layer 6] Audit Logging
    ↓
Response to User
```

## Configuration

### Environment Variables

All settings can be configured via `.env` (see [.env.example](agent-demo/.env.example)):

```bash
# API Configuration
CLAUDE_API_KEY=your_key_here
DEFAULT_MODEL=claude-sonnet-4-5-20250929
MAX_TOKENS=4096

# Memory Settings
CONVERSATION_HISTORY_LIMIT=50

# Observability
ENABLE_TRACING=true
LOG_LEVEL=INFO

# Agent Behavior
MAX_AGENT_ITERATIONS=10

# Governance
ENABLE_GOVERNANCE=true
GOVERNANCE_CONFIG=config/governance.yaml
```

### Governance Configuration

Production settings in [config/governance.yaml](agent-demo/config/governance.yaml):

```yaml
governance:
  enabled: true

  # NeMo Guardrails
  nemo_enabled: true
  nemo_config_path: "src/governance/rails_config"

  # PII Protection
  pii_enabled: true
  pii_language: "en"

  # Budget Controls
  budgets_enabled: true
  budgets_enforce: true
  global_daily_tokens: 1000000      # 1M tokens/day
  per_session_tokens: 50000         # 50K tokens/session
  per_user_daily_tokens: 100000     # 100K tokens/user/day

  # Audit
  audit_enabled: true
  audit_log_all_requests: true
  audit_retention_days: 90
```

For development, use relaxed settings:
```bash
export GOVERNANCE_CONFIG=config/governance.dev.yaml
```

## Development

### Adding a New Agent

1. Create a new file in `src/agents/`
2. Inherit from `BaseAgent`
3. Implement the `process()` method
4. Register with the coordinator

Example:
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

1. Define the tool schema in `src/tools/schemas.py`
2. Implement the tool function in `src/tools/registry.py`
3. Register it in the `_TOOL_REGISTRY` dict

Example:
```python
# 1. Define schema in tools/schemas.py
NEW_TOOL_SCHEMA = {
    "name": "new_tool",
    "description": "Tool description",
    "input_schema": {
        "type": "object",
        "properties": {...}
    }
}

# 2. Implement function in tools/registry.py
def new_tool(param1, param2):
    return result

# 3. Register
_TOOL_REGISTRY["new_tool"] = new_tool
```

### Running Tests

```bash
pytest tests/
```

## Mock Data

The system includes realistic mock data:
- **12 Artists** across various genres (Rock, Folk, Jazz, Electronic, etc.)
- **14 Venues** in Boston, MA and Nashville, TN
- Capacity ranges, genres, booking contacts, and pay ranges

Data is located in [data/mock_data.py](agent-demo/data/mock_data.py).

## Security & Governance

### Content Safety Features

- **Jailbreak Detection**: Blocks prompt injection attempts using NeMo Guardrails
- **Topic Control**: Keeps conversations focused on artist/venue matching
- **Toxic Content Filtering**: Prevents inappropriate language
- **Tool Validation**: Ensures agents only use authorized tools

### PII Protection

Presidio detects and protects 50+ PII types:
- Personal information (SSN, credit cards, names)
- Contact details (emails, phone numbers, addresses)
- Location data and IP addresses

Business contact information (artist/venue emails, venue phones) is allowed.

### Cost Controls

Budget hierarchy:
1. **Global daily limit**: 1M tokens/day across all users
2. **Per-session limit**: 50K tokens per session
3. **Per-user limit**: 100K tokens/user/day
4. **Rate limits**: 30 requests/min, 100K tokens/hour

### Audit Logging

All requests logged to `logs/audit/`:
- `requests.jsonl` - All requests (if enabled)
- `violations.jsonl` - Safety violations & PII events
- `metrics.jsonl` - Budget & performance metrics

## Performance

### Typical Metrics

**Latency:**
- Coordinator routing: ~3-4 seconds
- Tool execution: <100ms
- Specialist processing: ~3-10 seconds
- Governance overhead: 50-200ms per request
- Total end-to-end: ~7-17 seconds per query

**Token Usage:**
- Coordinator routing: ~850-950 tokens in per decision
- Specialist calls: ~900-1200 tokens in per call
- Tool usage: 70-150 tokens out per tool call
- Response generation: 300-500 tokens out

## Design Principles

### Clean Architecture
- Clear separation of concerns across 7 distinct layers
- Dependency injection throughout
- Abstract base classes for extensibility
- Single responsibility principle

### Observability First
- Every agent action is traced
- Structured logging with context
- Comprehensive metrics collection
- Real-time visualization

### Security by Design
- Defense-in-depth with 6 governance layers
- Agent-specific safety controls
- Complete audit trail
- Configurable enforcement modes

### Testability
- All components independently testable
- Mock data for deterministic testing
- Dependency injection enables easy mocking

## Technology Stack

- **LLM**: Claude Sonnet 4.5 via Anthropic API
- **Framework**: Custom multi-agent system
- **UI**: Streamlit
- **Configuration**: Pydantic Settings, python-dotenv
- **Observability**: Custom tracer, structured logging, metrics
- **Governance**:
  - NVIDIA NeMo Guardrails (agentic AI safety)
  - Microsoft Presidio (PII protection)
  - pyrate-limiter (cost control)

## Troubleshooting

### API errors?
- Check your `.env` file has a valid `CLAUDE_API_KEY`
- Verify the key at https://console.anthropic.com/settings/keys

### Import errors?
- The app automatically adds the project root to Python path
- Make sure you're in the `agent-demo` directory
- Run: `pip install -r requirements.txt`

### Presidio spaCy model missing?
```bash
python -m spacy download en_core_web_lg
```

### Budget limits too restrictive?
Adjust in `config/governance.yaml`:
```yaml
budgets_enforce: false  # Alert only
# OR increase limits:
per_session_tokens: 100000
```

### Disable governance entirely?
```bash
export ENABLE_GOVERNANCE=false
```

## Roadmap

### Phase 1: Foundation ✅
- Multi-agent orchestration
- Memory systems
- Observability infrastructure
- Production-ready governance

### Phase 2: Enhancement (Future)
- Persistent database storage
- User authentication
- Preference learning improvements
- Additional specialist agents
- API endpoint wrapper

### Phase 3: Production (Future)
- Real venue/artist data integration
- Deployment configuration
- Monitoring dashboards
- Distributed tracing (OpenTelemetry)
- Horizontal scaling

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## Support

For questions or issues, please open a GitHub issue.

---

Built with Claude Sonnet 4.5

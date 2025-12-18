# Artist-Venue Matching Multi-Agent System

A production-grade multi-agent system built with Claude that matches artists with venues using intelligent orchestration, memory management, and comprehensive observability.

## Features

### Multi-Agent Architecture
- **Coordinator Agent**: Intelligently routes requests to specialist agents
- **Artist Discovery Agent**: Searches and ranks artists by criteria
- **Venue Matching Agent**: Finds and scores venues for fit
- **Booking Advisor Agent**: Synthesizes information and provides recommendations

### Memory Management
- **Conversation Memory**: Maintains chat history across sessions
- **Working Memory**: Tracks intermediate results during multi-step tasks
- **Preference Memory**: Learns and persists user preferences

### Observability
- **Execution Tracing**: Detailed traces of all agent interactions
- **Structured Logging**: JSON-formatted logs with context
- **Performance Metrics**: Token usage and duration tracking per agent

### Interactive UI
- Real-time chat interface with Streamlit
- Split-view for chat and observability
- Live trace and metrics visualization

## Architecture

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
└── tests/               # Test suite

```

## Getting Started

### Prerequisites
- Python 3.10+
- Anthropic API key ([Get one here](https://console.anthropic.com/settings/keys))

### Installation

1. **Clone or navigate to the project**:
```bash
cd agent-demo
```

2. **Install dependencies**:
```bash
pip install -r requirements.txt
```

3. **Set up environment variables**:
```bash
cp .env.example .env
# Edit .env and add your CLAUDE_API_KEY
```

4. **Run the application**:
```bash
# Option 1: Direct run
streamlit run app/main.py

# Option 2: Using the launcher script
python run_app.py

# Option 3: Test first, then run
python test_system.py
streamlit run app/main.py
```

The application will open in your browser at `http://localhost:8501`

## Usage

### Example Queries

**Finding Venues:**
```
"I'm looking for rock venues in Boston with 200-500 capacity"
"Show me jazz clubs in Boston"
"What venues in Nashville book country music?"
```

**Finding Artists:**
```
"Find me some folk artists in Nashville"
"Show me punk bands in Boston"
"What electronic artists play venues under 1000 capacity?"
```

**Getting Recommendations:**
```
"Recommend artist-venue pairings for indie rock in Boston"
"I need to book an acoustic act for a 100-person venue"
"Match me some R&B artists with appropriate Boston venues"
```

### Observability Features

Toggle **"Show Traces & Metrics"** in the sidebar to see:

- **Metrics Panel**: Token usage, request counts, average durations by agent
- **Trace Viewer**: Step-by-step execution traces showing:
  - Agent routing decisions
  - Tool calls and results
  - LLM request/response timing
  - Memory operations

## Configuration

All settings can be configured via environment variables (see [.env.example](agent-demo/.env.example)):

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
```

## Development

### Project Structure

- **`src/agents/`**: All agent implementations inherit from `BaseAgent`
- **`src/memory/`**: Three memory systems (conversation, working, preferences)
- **`src/observability/`**: Tracer, logger, and metrics collector
- **`src/orchestration/`**: Main `AgentOrchestrator` that coordinates everything
- **`app/`**: Streamlit UI with modular components

### Adding a New Agent

1. Create a new file in `src/agents/`
2. Inherit from `BaseAgent`
3. Implement the `process()` method
4. Register with the coordinator in `orchestration/executor.py`

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

### Running Tests

```bash
pytest tests/
```

## Mock Data

The system includes mock data for:
- **12 Artists** across various genres (Rock, Folk, Jazz, Electronic, etc.)
- **14 Venues** in Boston, MA and Nashville, TN
- Realistic capacity ranges, genres, booking contacts, and pay ranges

Data is located in [data/mock_data.py](agent-demo/data/mock_data.py).

## Design Principles

### Clean Architecture
- Clear separation of concerns
- Dependency injection throughout
- Abstract base classes for extensibility

### Observability First
- Every agent action is traced
- Structured logging with context
- Comprehensive metrics collection

### Memory Management
- Conversation history maintained across sessions
- Working memory for multi-step tasks
- Preference learning (foundation for personalization)

### Testability
- All components independently testable
- Mock data for deterministic testing
- Dependency injection enables easy mocking

## Technical Stack

- **LLM**: Claude Sonnet 4.5 via Anthropic API
- **Framework**: Custom multi-agent system
- **UI**: Streamlit
- **Configuration**: Pydantic Settings
- **Observability**: Custom tracer, structured logging, metrics

## Roadmap

### Phase 1: Foundation ✅
- Multi-agent orchestration
- Memory systems
- Observability infrastructure

### Phase 2: Enhancement (Future)
- Persistent storage (database)
- User authentication
- Preference learning improvements
- Additional specialist agents

### Phase 3: Production (Future)
- Real venue/artist data integration
- API endpoints
- Deployment configuration
- Monitoring dashboards

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

Built with ❤️ using Claude Sonnet 4.5

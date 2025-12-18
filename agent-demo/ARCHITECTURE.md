# System Architecture

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Streamlit UI                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │     Chat     │  │    Traces    │  │   Metrics    │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Agent Orchestrator                           │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                  Coordinator Agent                        │  │
│  │              (Routes to specialists)                      │  │
│  └────────┬─────────────────────────────────────────────────┘  │
│           │                                                     │
│     ┌─────┴─────┬──────────────┬──────────────┐               │
│     ▼           ▼              ▼              ▼               │
│  ┌────────┐ ┌────────┐  ┌────────┐  ┌────────┐              │
│  │Artist  │ │ Venue  │  │Booking │  │Future  │              │
│  │Discovery│ │Matching│  │Advisor │  │Agents  │              │
│  └───┬────┘ └───┬────┘  └───┬────┘  └────────┘              │
│      │          │            │                                 │
└──────┼──────────┼────────────┼─────────────────────────────────┘
       │          │            │
       └──────────┴────────────┘
                  │
                  ▼
    ┌─────────────────────────────┐
    │       Tool Registry         │
    │  ┌──────────────────────┐   │
    │  │  search_artists      │   │
    │  │  search_venues       │   │
    │  │  get_artist_details  │   │
    │  │  get_venue_details   │   │
    │  └──────────────────────┘   │
    └─────────────┬───────────────┘
                  │
                  ▼
    ┌─────────────────────────────┐
    │        Mock Data            │
    │  12 Artists • 14 Venues     │
    └─────────────────────────────┘
```

## Component Layers

### Layer 1: User Interface
```
┌─────────────────────────────────────────┐
│          Streamlit Frontend             │
├─────────────────────────────────────────┤
│  • app/main.py                          │
│  • app/components/trace_viewer.py       │
│  • app/components/metrics_panel.py      │
└─────────────────────────────────────────┘
```

### Layer 2: Orchestration
```
┌─────────────────────────────────────────┐
│       Agent Orchestrator                │
├─────────────────────────────────────────┤
│  • Session management                   │
│  • Context building                     │
│  • Memory coordination                  │
│  • Error handling                       │
└─────────────────────────────────────────┘
```

### Layer 3: Agents
```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Coordinator  │  │   Artist     │  │    Venue     │
│    Agent     │─▶│  Discovery   │  │   Matching   │
│              │  │    Agent     │  │    Agent     │
└──────────────┘  └──────────────┘  └──────────────┘
                        │                   │
                        └───────┬───────────┘
                                │
                        ┌──────────────┐
                        │   Booking    │
                        │   Advisor    │
                        │    Agent     │
                        └──────────────┘
```

### Layer 4: Memory Systems
```
┌─────────────────────────────────────────┐
│            Memory Layer                 │
├─────────────────────────────────────────┤
│  ┌────────────────┐                     │
│  │ Conversation   │  Chat history       │
│  │    Memory      │  (50 msgs default)  │
│  └────────────────┘                     │
│  ┌────────────────┐                     │
│  │   Working      │  Task context       │
│  │    Memory      │  (request scope)    │
│  └────────────────┘                     │
│  ┌────────────────┐                     │
│  │  Preference    │  User preferences   │
│  │    Memory      │  (persistent)       │
│  └────────────────┘                     │
└─────────────────────────────────────────┘
```

### Layer 5: Observability
```
┌─────────────────────────────────────────┐
│        Observability Layer              │
├─────────────────────────────────────────┤
│  ┌────────────┐  ┌────────────┐        │
│  │   Tracer   │  │   Logger   │        │
│  │  (Events)  │  │   (JSON)   │        │
│  └────────────┘  └────────────┘        │
│  ┌────────────┐                         │
│  │  Metrics   │  Token usage, timing   │
│  │ Collector  │  per agent             │
│  └────────────┘                         │
└─────────────────────────────────────────┘
```

### Layer 6: Tools & Data
```
┌─────────────────────────────────────────┐
│           Tool Layer                    │
├─────────────────────────────────────────┤
│  • Tool schemas (JSON)                  │
│  • Tool registry (functions)            │
│  • Tool execution dispatcher            │
└─────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│           Data Layer                    │
├─────────────────────────────────────────┤
│  • Mock artists data                    │
│  • Mock venues data                     │
│  • (Extensible to real APIs)            │
└─────────────────────────────────────────┘
```

## Request Flow

### 1. User Query Processing
```
User Input
    │
    ▼
Streamlit UI receives input
    │
    ▼
AgentOrchestrator.process_message()
    │
    ├─▶ Add to conversation memory
    │
    ├─▶ Build context (preferences, etc.)
    │
    ├─▶ Get conversation history
    │
    └─▶ Start trace
```

### 2. Agent Routing
```
Coordinator Agent receives query
    │
    ▼
Analyze intent
    │
    ├─▶ Artist search? → Artist Discovery Agent
    │
    ├─▶ Venue search? → Venue Matching Agent
    │
    ├─▶ Booking advice? → Booking Advisor Agent
    │
    └─▶ General query? → Direct response
```

### 3. Specialist Processing
```
Specialist Agent
    │
    ▼
Build messages with context
    │
    ▼
Call LLM (with tools)
    │
    ├─▶ Tool use? → Execute tool(s)
    │     │
    │     └─▶ Add results to messages
    │           │
    │           └─▶ Loop back to LLM
    │
    └─▶ Text response → Return to coordinator
```

### 4. Response & Observability
```
Response flows back
    │
    ├─▶ Add to conversation memory
    │
    ├─▶ Record metrics (tokens, timing)
    │
    ├─▶ Complete trace
    │
    └─▶ Return to UI
        │
        └─▶ Display to user
            │
            └─▶ Show observability (optional)
```

## Data Flow Diagram

```
┌─────────┐
│  User   │
└────┬────┘
     │ Query
     ▼
┌─────────────────┐      ┌─────────────────┐
│   Streamlit     │─────▶│  Orchestrator   │
│      UI         │◀─────│                 │
└─────────────────┘      └────────┬────────┘
     │                            │
     │                            ▼
     │                   ┌─────────────────┐
     │                   │  Coordinator    │
     │                   └────────┬────────┘
     │                            │
     │              ┌─────────────┼─────────────┐
     │              ▼             ▼             ▼
     │        ┌──────────┐  ┌──────────┐  ┌──────────┐
     │        │ Artist   │  │  Venue   │  │ Booking  │
     │        │Discovery │  │ Matching │  │ Advisor  │
     │        └────┬─────┘  └────┬─────┘  └────┬─────┘
     │             │             │             │
     │             └──────┬──────┴──────┬──────┘
     │                    │             │
     │                    ▼             │
     │             ┌─────────────┐      │
     │             │   Tools     │      │
     │             └─────────────┘      │
     │                    │             │
     │                    ▼             │
     │             ┌─────────────┐      │
     │             │    Data     │      │
     │             └─────────────┘      │
     │                                  │
     └──────────────────────────────────┘
           Response + Observability
```

## Memory Architecture

```
┌────────────────────────────────────────────────┐
│              Memory Systems                    │
├────────────────────────────────────────────────┤
│                                                │
│  ┌──────────────────────────────────────┐     │
│  │     Conversation Memory              │     │
│  │  ┌────────────────────────────────┐  │     │
│  │  │  Session 1: [msg1, msg2, ...]  │  │     │
│  │  │  Session 2: [msg1, msg2, ...]  │  │     │
│  │  │  Session N: [msg1, msg2, ...]  │  │     │
│  │  └────────────────────────────────┘  │     │
│  │  Lifecycle: Until explicitly cleared │     │
│  └──────────────────────────────────────┘     │
│                                                │
│  ┌──────────────────────────────────────┐     │
│  │     Working Memory                   │     │
│  │  ┌────────────────────────────────┐  │     │
│  │  │  Context ID: "abc123"          │  │     │
│  │  │  ├─ User Query                 │  │     │
│  │  │  ├─ Intent: "venue_search"     │  │     │
│  │  │  ├─ Intermediate Results       │  │     │
│  │  │  └─ Routing Decisions          │  │     │
│  │  └────────────────────────────────┘  │     │
│  │  Lifecycle: Single request           │     │
│  └──────────────────────────────────────┘     │
│                                                │
│  ┌──────────────────────────────────────┐     │
│  │     Preference Memory                │     │
│  │  ┌────────────────────────────────┐  │     │
│  │  │  User 1:                       │  │     │
│  │  │  ├─ Genres: [Rock, Jazz]       │  │     │
│  │  │  ├─ Locations: [Boston]        │  │     │
│  │  │  └─ Capacity: (100-500)        │  │     │
│  │  └────────────────────────────────┘  │     │
│  │  Lifecycle: Persistent per user      │     │
│  └──────────────────────────────────────┘     │
│                                                │
└────────────────────────────────────────────────┘
```

## Observability Architecture

```
┌────────────────────────────────────────────────┐
│         Observability Stack                    │
├────────────────────────────────────────────────┤
│                                                │
│  Every Agent Action Generates:                │
│                                                │
│  ┌──────────────────────────────────────┐     │
│  │         Trace Event                  │     │
│  │  • Event Type                        │     │
│  │  • Agent Name                        │     │
│  │  • Timestamp                         │     │
│  │  • Duration (if applicable)          │     │
│  │  • Event Data                        │     │
│  └──────────────────────────────────────┘     │
│           │                                    │
│           ├──▶ Stored in Trace                │
│           │                                    │
│           ├──▶ Logged (JSON format)           │
│           │                                    │
│           └──▶ Metrics updated                │
│                                                │
│  End Result:                                   │
│  • Complete execution trace                   │
│  • Structured logs                            │
│  • Performance metrics                        │
│  • Token usage tracking                       │
│                                                │
└────────────────────────────────────────────────┘
```

## Tool Execution Flow

```
Agent needs information
        │
        ▼
Determines which tool to call
        │
        ▼
LLM generates tool_use block
        │
        ├─ tool_name
        ├─ tool_id
        └─ input parameters
        │
        ▼
Agent extracts tool call
        │
        ▼
Tool Registry dispatcher
        │
        ├─▶ Lookup tool function
        ├─▶ Execute with parameters
        └─▶ Return result
        │
        ▼
Format as tool_result
        │
        ▼
Add to message history
        │
        ▼
Continue agent loop
```

## Configuration Flow

```
.env file
    │
    ├─ CLAUDE_API_KEY
    ├─ DEFAULT_MODEL
    ├─ MAX_TOKENS
    ├─ LOG_LEVEL
    └─ ... (other settings)
    │
    ▼
Pydantic Settings
    │
    ├─ Validation
    ├─ Type coercion
    └─ Default values
    │
    ▼
Global settings object
    │
    ├─▶ Used by Agents
    ├─▶ Used by Orchestrator
    ├─▶ Used by Logger
    └─▶ Used by Tracer
```

## Extension Points

```
┌─────────────────────────────────────────┐
│      Where to Add New Features          │
├─────────────────────────────────────────┤
│                                         │
│  New Agent:                             │
│  └─▶ src/agents/your_agent.py          │
│      └─▶ Register in orchestrator      │
│                                         │
│  New Tool:                              │
│  ├─▶ Schema: src/tools/schemas.py      │
│  └─▶ Impl: src/tools/registry.py       │
│                                         │
│  New Memory Store:                      │
│  └─▶ src/memory/your_memory.py         │
│      └─▶ Add to orchestrator           │
│                                         │
│  New UI Component:                      │
│  └─▶ app/components/your_comp.py       │
│      └─▶ Import in main.py             │
│                                         │
│  New Configuration:                     │
│  ├─▶ Add to .env.example               │
│  └─▶ Add to src/config/settings.py     │
│                                         │
└─────────────────────────────────────────┘
```

## Technology Stack

```
┌─────────────────────────────────────────┐
│         Technology Layers               │
├─────────────────────────────────────────┤
│  Frontend:                              │
│  └─▶ Streamlit 1.31+                    │
│                                         │
│  LLM:                                   │
│  └─▶ Claude Sonnet 4.5                  │
│                                         │
│  Framework:                             │
│  ├─▶ Custom multi-agent system          │
│  └─▶ Anthropic Python SDK               │
│                                         │
│  Configuration:                         │
│  ├─▶ Pydantic 2.0+                      │
│  └─▶ python-dotenv                      │
│                                         │
│  Data:                                  │
│  └─▶ Pandas (for metrics viz)           │
│                                         │
│  Dev Tools:                             │
│  ├─▶ pytest                             │
│  ├─▶ black                              │
│  ├─▶ ruff                               │
│  └─▶ mypy                               │
└─────────────────────────────────────────┘
```

This architecture enables:
- ✅ Scalable agent addition
- ✅ Pluggable components
- ✅ Full observability
- ✅ Easy testing
- ✅ Production readiness

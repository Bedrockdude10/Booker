# Quick Start Guide

Get the multi-agent system running in 3 minutes!

## Prerequisites

- Python 3.10+
- Your Anthropic API key

## Installation & Setup

```bash
# 1. Navigate to the project
cd agent-demo

# 2. Install dependencies
pip install -r requirements.txt

# 3. The .env file is already configured with your API key
# (If you need to change it, edit .env)

# 4. Test the system (optional but recommended)
python test_system.py

# 5. Launch the Streamlit app
streamlit run app/main.py
```

The app will open at `http://localhost:8501`

## First Steps in the UI

1. **Try a simple query**: "Find me rock venues in Boston"
   - Watch the coordinator route it to the venue matching agent
   - See the results with explanations

2. **Toggle observability**: Click "Show Traces & Metrics" in the sidebar
   - View execution traces in real-time
   - See token usage by agent
   - Explore routing decisions

3. **Try different queries**:
   - "Show me folk artists in Nashville"
   - "I need an indie rock venue with 200-300 capacity"
   - "Recommend some artist-venue pairings for jazz"

## What to Look For

### Multi-Agent Routing
- The coordinator analyzes your query
- Routes to the appropriate specialist
- You'll see "🧭 Routed to: [agent_name]" below responses

### Tool Usage
- Agents call search tools to find matches
- Tool calls are visible in the trace viewer
- Results are synthesized into helpful responses

### Memory
- Your conversation history is maintained
- Try referencing previous queries
- Session metrics track all activity

## Troubleshooting

**Import errors?**
- The app automatically adds the project root to Python path
- If issues persist, run from the project root directory

**API errors?**
- Check your `.env` file has a valid `CLAUDE_API_KEY`
- Verify the key at https://console.anthropic.com/settings/keys

**Module not found?**
- Make sure you're in the `agent-demo` directory
- Run: `pip install -r requirements.txt`

## Next Steps

- Read [README.md](README.md) for full documentation
- Explore the codebase in `src/`
- Add new agents or tools
- Customize prompts in `src/config/prompts.py`

## System Architecture Overview

```
User Query → Coordinator Agent → Specialist Agents → Tools → Response
                 ↓                      ↓               ↓
             Tracer             Memory Systems     Metrics
```

**4 Agents:**
- 🧭 Coordinator (router)
- 🎸 Artist Discovery
- 🏛️ Venue Matching
- 💡 Booking Advisor

**3 Memory Systems:**
- 💬 Conversation
- 🧠 Working Memory
- ⭐ Preferences

**Full Observability:**
- 🔍 Traces
- 📊 Metrics
- 📝 Logs

Enjoy exploring the multi-agent system! 🎵

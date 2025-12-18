"""
Main Streamlit application for the multi-agent artist-venue matching system.
"""

import sys
from pathlib import Path

# Add project root to path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

import streamlit as st
import uuid
from datetime import datetime

from src.orchestration.executor import AgentOrchestrator
from app.components.trace_viewer import display_trace_summary
from app.components.metrics_panel import display_session_metrics, display_token_breakdown

# Page configuration
st.set_page_config(
    page_title="Artist-Venue Matching Agent",
    page_icon="🎵",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Initialize session state
if "session_id" not in st.session_state:
    st.session_state.session_id = str(uuid.uuid4())

if "messages" not in st.session_state:
    st.session_state.messages = []

if "orchestrator" not in st.session_state:
    st.session_state.orchestrator = AgentOrchestrator()

if "show_observability" not in st.session_state:
    st.session_state.show_observability = False

# Main title
st.title("🎵 Artist-Venue Matching Agent")
st.caption("Powered by Claude Multi-Agent System • Connecting artists with venues")

# Sidebar
with st.sidebar:
    st.header("About")
    st.markdown("""
    This **multi-agent system** helps you:
    - 🎸 Find venues for your band
    - 🎤 Discover artists for your venue
    - 🔍 Get booking recommendations

    **Example queries:**
    - "I'm looking for rock venues in Boston with 200-500 capacity"
    - "Find me some folk artists in Nashville"
    - "What venues in Boston book electronic music?"
    - "Recommend some artist-venue pairings for indie rock"
    """)

    st.divider()

    st.header("System Architecture")
    st.markdown("""
    **Agents:**
    - 🧭 **Coordinator**: Routes requests
    - 🎸 **Artist Discovery**: Finds artists
    - 🏛️ **Venue Matching**: Finds venues
    - 💡 **Booking Advisor**: Provides recommendations

    **Features:**
    - 🔍 Multi-agent orchestration
    - 💾 Conversation memory
    - 📊 Real-time observability
    - 🎯 Smart routing
    """)

    st.divider()

    # Observability toggle
    st.header("Observability")
    st.session_state.show_observability = st.toggle(
        "Show Traces & Metrics",
        value=st.session_state.show_observability
    )

    st.divider()

    # Session controls
    st.header("Session")
    st.caption(f"ID: {st.session_state.session_id[:8]}...")

    if st.button("🗑️ Clear Chat History", use_container_width=True):
        st.session_state.orchestrator.clear_session(st.session_state.session_id)
        st.session_state.messages = []
        st.rerun()

    st.caption(f"Messages: {len(st.session_state.messages)}")

# Main content area
if st.session_state.show_observability:
    # Split layout for chat and observability
    chat_col, obs_col = st.columns([1, 1])

    with chat_col:
        st.subheader("💬 Chat")
        chat_container = st.container(height=500)

    with obs_col:
        st.subheader("📊 Observability")
        obs_container = st.container(height=500)
else:
    # Full-width chat
    chat_container = st.container()
    obs_container = None

# Display chat messages
with chat_container:
    for message in st.session_state.messages:
        with st.chat_message(message["role"]):
            st.markdown(message["content"])

            # Show metadata if available
            if message.get("metadata"):
                metadata = message["metadata"]
                if metadata.get("routed"):
                    st.caption(f"🧭 Routed to: {metadata.get('target_agent', 'unknown')}")

# Display observability if enabled
if st.session_state.show_observability and obs_container:
    with obs_container:
        tab1, tab2 = st.tabs(["📊 Metrics", "🔍 Traces"])

        with tab1:
            metrics = st.session_state.orchestrator.get_session_metrics(st.session_state.session_id)
            display_session_metrics(metrics)

            if metrics:
                st.divider()
                display_token_breakdown(metrics)

        with tab2:
            traces = st.session_state.orchestrator.get_recent_traces(limit=5)
            display_trace_summary(traces)

# Chat input
if user_input := st.chat_input("Ask me to find venues or artists..."):
    # Add user message to chat
    st.session_state.messages.append({"role": "user", "content": user_input})

    # Display user message immediately
    with chat_container:
        with st.chat_message("user"):
            st.markdown(user_input)

    # Process through orchestrator
    with chat_container:
        with st.chat_message("assistant"):
            with st.spinner("Thinking..."):
                result = st.session_state.orchestrator.process_message(
                    user_input,
                    st.session_state.session_id,
                    user_id=None  # Could add user authentication here
                )

            # Display response
            st.markdown(result["content"])

            # Show metadata
            metadata = result.get("metadata", {})
            if metadata.get("routed"):
                st.caption(f"🧭 Routed to: {metadata.get('target_agent', 'unknown')}")

            # Show token usage
            tokens = result.get("tokens", {})
            if tokens.get("total", 0) > 0:
                st.caption(f"💬 Tokens: {tokens['in']} in, {tokens['out']} out ({tokens['total']} total)")

    # Add assistant message to chat history
    st.session_state.messages.append({
        "role": "assistant",
        "content": result["content"],
        "metadata": metadata
    })

    # Rerun to update observability panel
    if st.session_state.show_observability:
        st.rerun()

# Footer
st.divider()
st.caption("Multi-Agent System • Built with Claude Sonnet 4.5 • Mock data for Boston, MA and Nashville, TN")

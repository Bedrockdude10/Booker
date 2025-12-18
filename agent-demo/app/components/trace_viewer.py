"""Trace viewer component for displaying execution traces."""

import streamlit as st
from typing import Any


def display_trace(trace: dict[str, Any]) -> None:
    """Display an execution trace with all events.

    Args:
        trace: Trace dictionary from the tracer
    """
    with st.expander(f"🔍 Trace: {trace['trace_id'][:8]}... ({trace['event_count']} events)", expanded=False):
        # Summary metrics
        col1, col2, col3, col4 = st.columns(4)
        with col1:
            st.metric("Duration", f"{trace.get('duration_ms', 0):.0f}ms")
        with col2:
            st.metric("Tokens In", trace['total_tokens_in'])
        with col3:
            st.metric("Tokens Out", trace['total_tokens_out'])
        with col4:
            st.metric("Total Tokens", trace['total_tokens'])

        # Event timeline
        st.subheader("Event Timeline")
        for event in trace['events']:
            _display_event(event)


def _display_event(event: dict[str, Any]) -> None:
    """Display a single trace event.

    Args:
        event: Event dictionary
    """
    event_type = event['event_type']
    agent_name = event['agent_name']
    data = event['data']

    # Choose icon based on event type
    icon = _get_event_icon(event_type)

    # Format duration if present
    duration_str = ""
    if event.get('duration_ms'):
        duration_str = f" ({event['duration_ms']:.0f}ms)"

    # Display event
    with st.container():
        st.markdown(f"{icon} **{event_type}** - {agent_name}{duration_str}")

        # Show event data in a compact format
        if data:
            if event_type == "tool_call":
                st.caption(f"🔧 {data.get('tool', 'unknown')}")
                if 'input' in data:
                    with st.expander("Tool Input", expanded=False):
                        st.json(data['input'])
            elif event_type == "routing_decision":
                st.caption(f"➡️ Routing to: {data.get('target_agent', 'unknown')}")
                if 'reason' in data:
                    st.caption(f"Reason: {data['reason']}")
            elif event_type in ["llm_request", "llm_response"]:
                st.caption(f"Tokens: {data.get('tokens_in', 0)} in, {data.get('tokens_out', 0)} out")
            else:
                with st.expander("Event Data", expanded=False):
                    st.json(data)


def _get_event_icon(event_type: str) -> str:
    """Get icon for event type.

    Args:
        event_type: Type of event

    Returns:
        Icon string
    """
    icons = {
        "agent_start": "🚀",
        "agent_end": "✅",
        "tool_call": "🔧",
        "tool_result": "📊",
        "llm_request": "💬",
        "llm_response": "🤖",
        "routing_decision": "🧭",
        "memory_read": "📖",
        "memory_write": "📝",
        "error": "❌"
    }
    return icons.get(event_type, "•")


def display_trace_summary(traces: list[dict[str, Any]]) -> None:
    """Display summary of multiple traces.

    Args:
        traces: List of trace dictionaries
    """
    if not traces:
        st.info("No traces available yet. Start a conversation to see traces.")
        return

    st.subheader("Recent Traces")

    for trace in traces:
        display_trace(trace)

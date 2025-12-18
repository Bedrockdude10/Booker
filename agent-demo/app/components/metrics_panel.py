"""Metrics panel component for displaying session metrics."""

import streamlit as st
from typing import Any


def display_session_metrics(metrics: dict[str, Any]) -> None:
    """Display metrics for the current session.

    Args:
        metrics: Metrics dictionary from the metrics collector
    """
    if not metrics:
        st.info("No metrics available yet. Start a conversation to see metrics.")
        return

    st.subheader("Session Metrics")

    # Overall session metrics
    st.markdown("### Overall")
    col1, col2, col3 = st.columns(3)

    with col1:
        st.metric("Total Requests", metrics.get('total_requests', 0))
    with col2:
        st.metric("Total Tokens", metrics.get('total_tokens', 0))
    with col3:
        tokens_in = metrics.get('total_tokens_in', 0)
        tokens_out = metrics.get('total_tokens_out', 0)
        if tokens_in + tokens_out > 0:
            ratio = tokens_out / (tokens_in + tokens_out) * 100
            st.metric("Output Ratio", f"{ratio:.1f}%")

    # Agent-specific metrics
    agents = metrics.get('agents', {})
    if agents:
        st.markdown("### By Agent")

        for agent_name, agent_metrics in agents.items():
            with st.expander(f"🤖 {agent_name}", expanded=True):
                col1, col2, col3, col4 = st.columns(4)

                with col1:
                    st.metric("Calls", agent_metrics.get('total_calls', 0))
                with col2:
                    st.metric("Avg Duration", f"{agent_metrics.get('avg_duration_ms', 0):.0f}ms")
                with col3:
                    st.metric("Total Tokens", agent_metrics.get('total_tokens', 0))
                with col4:
                    st.metric("Errors", agent_metrics.get('error_count', 0))


def display_token_breakdown(metrics: dict[str, Any]) -> None:
    """Display token usage breakdown by agent.

    Args:
        metrics: Metrics dictionary from the metrics collector
    """
    agents = metrics.get('agents', {})
    if not agents:
        return

    st.subheader("Token Distribution")

    # Create data for visualization
    agent_names = []
    token_counts = []

    for agent_name, agent_metrics in agents.items():
        agent_names.append(agent_name)
        token_counts.append(agent_metrics.get('total_tokens', 0))

    # Display as bar chart
    if token_counts and sum(token_counts) > 0:
        import pandas as pd

        df = pd.DataFrame({
            'Agent': agent_names,
            'Tokens': token_counts
        })

        st.bar_chart(df.set_index('Agent'))
    else:
        st.info("No token data available yet.")

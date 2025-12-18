"""
Quick test script to verify the multi-agent system works.
"""

import sys
from pathlib import Path

# Add project root to path
project_root = Path(__file__).parent
sys.path.insert(0, str(project_root))

from src.orchestration.executor import AgentOrchestrator


def test_basic_query():
    """Test a basic query through the system."""
    print("🧪 Testing Multi-Agent System\n")
    print("=" * 60)

    # Initialize orchestrator
    print("\n1. Initializing orchestrator...")
    orchestrator = AgentOrchestrator()
    print("   ✓ Orchestrator initialized")

    # Test query
    query = "Find me rock venues in Boston"
    session_id = "test_session"

    print(f"\n2. Processing query: '{query}'")
    print("   (This may take a few moments...)")

    try:
        result = orchestrator.process_message(query, session_id)

        print("\n3. ✅ Success!")
        print("\n" + "=" * 60)
        print("RESPONSE:")
        print("=" * 60)
        print(result["content"])
        print("\n" + "=" * 60)
        print("METADATA:")
        print("=" * 60)
        print(f"   Trace ID: {result['trace_id']}")
        print(f"   Tokens: {result['tokens']}")
        if result.get("metadata"):
            print(f"   Routed to: {result['metadata'].get('target_agent', 'N/A')}")
        print("=" * 60)

        # Get metrics
        metrics = orchestrator.get_session_metrics(session_id)
        if metrics:
            print("\nSESSION METRICS:")
            print("=" * 60)
            print(f"   Total Requests: {metrics.get('total_requests', 0)}")
            print(f"   Total Tokens: {metrics.get('total_tokens', 0)}")
            print("=" * 60)

        print("\n✅ All tests passed! System is working correctly.")
        print("\nYou can now run: streamlit run app/main.py")

    except Exception as e:
        print(f"\n❌ Error: {str(e)}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    test_basic_query()

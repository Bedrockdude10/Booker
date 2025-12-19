# Agentic AI Governance Layer

Production-ready safety & governance for the Booker multi-agent system.

## Overview

This governance layer provides comprehensive safety controls specifically designed for **agentic AI systems**, using industry-leading libraries:

- **🎯 NVIDIA NeMo Guardrails** - Complete agentic AI safety framework
- **🔒 Microsoft Presidio** - Industry-standard PII detection & anonymization
- **💰 pyrate-limiter** - Cost control & rate limiting

## Why This Stack?

### Agent-Focused Design

Unlike generic LLM safety tools, this stack is purpose-built for **multi-agent systems**:

- **NeMo Guardrails**: Validates multi-step agent reasoning, tool calling, and dialogue flows
- **Agent-aware topic control**: Keeps agents focused on artist/venue matching tasks
- **Tool usage validation**: Ensures agents only use authorized tools
- **70% less complexity**: 3 libraries instead of 7, all agent-optimized

### Key Features

✅ **Jailbreak Detection** - Blocks prompt injection attempts
✅ **Content Safety** - Filters toxic/inappropriate content
✅ **Topic Control** - Prevents conversation drift
✅ **PII Protection** - Detects & anonymizes personal information
✅ **Cost Controls** - Budget limits & rate limiting
✅ **Audit Logging** - Complete compliance trail
✅ **Tool Safety** - Validates agent tool usage

## Architecture

### Defense-in-Depth (6 Layers)

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
  ↓
[Layer 5] Output Filtering & PII Protection
  ↓
[Layer 6] Audit Logging
  ↓
Response to User
```

### Components

```
src/governance/
├── coordinator.py          # Main governance orchestrator
├── nemo_wrapper.py         # NeMo Guardrails integration
├── pii_protection.py       # Presidio PII detection
├── cost_control.py         # pyrate-limiter budget control
├── audit.py                # Compliance audit logger
├── models.py               # Data models
└── rails_config/           # NeMo Guardrails configuration
    └── config.yml          # Colang safety rules
```

## Installation

### 1. Install Dependencies

```bash
cd agent-demo
pip install -r requirements.txt
```

This installs:
- `nemoguardrails>=0.9.0` - Agentic AI safety
- `presidio-analyzer>=2.2.0` - PII detection
- `presidio-anonymizer>=2.2.0` - PII anonymization
- `pyrate-limiter>=3.0.0` - Rate limiting
- `pyyaml>=6.0.0` - Config loading

### 2. Download spaCy Model (for Presidio)

```bash
python -m spacy download en_core_web_lg
```

### 3. Configure Environment

Add to `.env`:

```bash
# Governance (optional - defaults shown)
ENABLE_GOVERNANCE=true
GOVERNANCE_CONFIG=config/governance.yaml
```

## Configuration

### Production Config (`config/governance.yaml`)

```yaml
governance:
  enabled: true

  # NeMo Guardrails
  nemo_enabled: true
  nemo_config_path: "src/governance/rails_config"

  # PII Protection
  pii_enabled: true
  pii_language: "en"
  pii_allowed_entities:
    - "ARTIST_EMAIL"    # Allow business contact info
    - "VENUE_PHONE"
    - "VENUE_EMAIL"

  # Budget Controls
  budgets_enabled: true
  budgets_enforce: true
  global_daily_tokens: 1000000      # 1M tokens/day
  per_session_tokens: 50000         # 50K tokens/session
  per_user_daily_tokens: 100000     # 100K tokens/user/day
  requests_per_minute: 30
  tokens_per_hour: 100000

  # Audit
  audit_enabled: true
  audit_log_all_requests: true
  audit_retention_days: 90
```

### Development Config (`config/governance.dev.yaml`)

Use for local development with relaxed settings:

```bash
export GOVERNANCE_CONFIG=config/governance.dev.yaml
```

## Usage

### Basic Integration

```python
from src.config.settings import settings
from src.governance import GovernanceCoordinator

# Load governance config
gov_config = settings.load_governance_config()

# Initialize coordinator
governance = GovernanceCoordinator(gov_config)

# Check user input
result = governance.check_input(
    user_input="Find me a venue in Boston",
    session_id="session-123",
    estimated_tokens=1000
)

if not result["passed"]:
    print(f"Request blocked: {result['violations']}")
else:
    # Use sanitized input
    safe_input = result["sanitized_input"]

    # Process with agents...
    # ...

    # Check agent output
    output_result = governance.check_output(
        agent_output=response,
        session_id="session-123",
        event_id=result["event_id"]
    )

    # Commit token usage
    governance.commit_usage(
        reservation_id=result["reservation_id"],
        actual_tokens=tokens_used,
        session_id="session-123"
    )
```

### Get Statistics

```python
stats = governance.get_stats()
print(stats)
# {
#   "enabled": True,
#   "nemo": {"enabled": True, "rails_loaded": True},
#   "pii_protector": {"enabled": True, "language": "en"},
#   "cost_controller": {"global_usage": 45000, ...},
#   "audit_logger": {"requests_count": 127, ...}
# }
```

## NeMo Guardrails Configuration

### Allowed Topics

Defined in `src/governance/rails_config/config.yml`:

```yaml
define valid booking topics:
  "artists"
  "venues"
  "concerts"
  "gigs"
  "bookings"
```

### Blocked Topics

```yaml
define user ask off topic:
  "politics"
  "violence"
  "illegal activities"
```

### Allowed Tools

```yaml
define allowed tools:
  "search_artists"
  "search_venues"
  "get_artist_details"
  "get_venue_details"
```

## PII Protection

### Detected Entity Types

Presidio detects 50+ PII types including:
- EMAIL_ADDRESS
- PHONE_NUMBER
- CREDIT_CARD
- SSN
- IP_ADDRESS
- PERSON (names)
- LOCATION

### Business Contact Allowlist

Artist/venue contact information is **allowed** (business PII):

```python
pii_allowed_entities = [
    "ARTIST_EMAIL",
    "VENUE_PHONE",
    "VENUE_EMAIL"
]
```

User personal information is **protected**.

## Cost Control

### Budget Hierarchy

1. **Global daily limit**: 1M tokens/day across all users
2. **Per-session limit**: 50K tokens per session
3. **Per-user limit**: 100K tokens/user/day
4. **Rate limits**: 30 requests/min, 100K tokens/hour

### Enforcement Modes

**Production** (`budgets_enforce: true`):
- Hard limits, requests blocked when exceeded

**Development** (`budgets_enforce: false`):
- Alert only, requests allowed but logged

## Audit Logging

### Log Files

Audit logs are written to `logs/audit/`:

```
logs/audit/
├── requests.jsonl       # All requests (if enabled)
├── violations.jsonl     # Safety violations & PII events
└── metrics.jsonl        # Budget & performance metrics
```

### Log Format

JSON Lines (`.jsonl`) for easy parsing:

```json
{
  "event_id": "550e8400-e29b-41d4-a716-446655440000",
  "event_type": "request",
  "timestamp": "2025-12-19T10:30:00Z",
  "session_id": "session-123",
  "user_id": "user-456",
  "data": {"input": "Find artists in Nashville"},
  "guardrails_passed": true,
  "pii_detected": false,
  "budget_allowed": true
}
```

## Performance

### Latency Overhead

Expected: **50-200ms per request**

- Input validation: 5-10ms
- NeMo Guardrails: 20-100ms (can be cached)
- Presidio PII: 10-40ms
- pyrate-limiter: <5ms
- Audit logging: <10ms (async)

### Optimization

- Model caching (Detoxify, Presidio)
- Async audit writes
- Conditional checks (dev vs prod)
- Redis support for distributed rate limiting

## Monitoring

### Health Checks

```python
# Check governance health
stats = governance.get_stats()

if not stats["enabled"]:
    print("⚠️  Governance layer is disabled!")

if not stats["nemo"]["rails_loaded"]:
    print("⚠️  NeMo Guardrails failed to load!")
```

### Key Metrics

Monitor these in production:

- Budget usage vs limits
- Violation rate
- PII detection frequency
- Request latency (p50, p95, p99)
- Guardrail failure rate

## Troubleshooting

### NeMo Guardrails not loading

```
Error: Failed to initialize NeMo Guardrails
```

**Solution**: Check that `src/governance/rails_config/config.yml` exists and is valid YAML.

### Presidio spaCy model missing

```
Error: Can't find model 'en_core_web_lg'
```

**Solution**:
```bash
python -m spacy download en_core_web_lg
```

### Budget limits too restrictive

**Solution**: Adjust in `config/governance.yaml`:

```yaml
budgets_enforce: false  # Alert only
# OR increase limits:
per_session_tokens: 100000
```

## Rollback / Disable

### Disable All Governance

```bash
export ENABLE_GOVERNANCE=false
```

### Disable Individual Components

In `config/governance.yaml`:

```yaml
governance:
  enabled: true
  nemo_enabled: false          # Disable NeMo
  pii_enabled: false           # Disable PII protection
  budgets_enabled: false       # Disable cost controls
  audit_enabled: false         # Disable audit logging
```

### Shadow Mode (Monitoring Only)

```yaml
budgets_enforce: false         # Log violations but don't block
audit_log_all_requests: true   # Log everything for analysis
```

## Security Best Practices

1. **Always enable in production**: Set `ENABLE_GOVERNANCE=true`
2. **Monitor audit logs**: Review `logs/audit/violations.jsonl` regularly
3. **Tune thresholds**: Adjust based on actual usage patterns
4. **Protect API keys**: Never commit `.env` files
5. **Regular updates**: Keep NeMo, Presidio, and dependencies updated

## Next Steps

### Integration with Orchestrator

To integrate with the agent orchestrator:

1. Initialize governance in orchestrator
2. Wrap `process_message()` with governance checks
3. Use sanitized inputs for agent processing
4. Validate outputs before returning to user
5. Commit token usage after completion

See `src/orchestration/executor.py` for integration points.

### Custom Rails

Add custom safety rules in `src/governance/rails_config/config.yml`:

```yaml
define user ask for price:
  "how much"
  "cost"
  "price"

define bot explain pricing:
  "Pricing varies by venue. Please contact the venue directly for rates."
```

## Support

- **NeMo Guardrails**: https://github.com/NVIDIA-NeMo/Guardrails
- **Presidio**: https://github.com/microsoft/presidio
- **pyrate-limiter**: https://github.com/vutran1710/PyrateLimiter

## License

Same as parent project.

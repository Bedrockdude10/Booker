"""Application settings from environment variables."""

import os
import yaml
from pathlib import Path
from pydantic_settings import BaseSettings
from pydantic import Field

from src.governance.models import GovernanceConfig


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # Anthropic API
    anthropic_api_key: str = Field(..., alias="CLAUDE_API_KEY")
    default_model: str = Field("claude-sonnet-4-5-20250929", alias="DEFAULT_MODEL")
    max_tokens: int = Field(4096, alias="MAX_TOKENS")

    # Memory
    conversation_history_limit: int = Field(50, alias="CONVERSATION_HISTORY_LIMIT")

    # Observability
    enable_tracing: bool = Field(True, alias="ENABLE_TRACING")
    log_level: str = Field("INFO", alias="LOG_LEVEL")

    # Agent settings
    max_agent_iterations: int = Field(10, alias="MAX_AGENT_ITERATIONS")

    # Governance settings
    enable_governance: bool = Field(True, alias="ENABLE_GOVERNANCE")
    governance_config_file: str = Field("config/governance.yaml", alias="GOVERNANCE_CONFIG")

    class Config:
        env_file = ".env"
        case_sensitive = False

    def load_governance_config(self) -> GovernanceConfig:
        """Load governance configuration from YAML file"""
        if not self.enable_governance:
            return GovernanceConfig(enabled=False)

        config_path = Path(self.governance_config_file)
        if not config_path.exists():
            print(f"Warning: Governance config not found at {config_path}, using defaults")
            return GovernanceConfig()

        try:
            with open(config_path) as f:
                yaml_config = yaml.safe_load(f)

            # Extract governance section
            gov_config = yaml_config.get("governance", {})

            return GovernanceConfig(**gov_config)
        except Exception as e:
            print(f"Error loading governance config: {e}")
            return GovernanceConfig()


# Global settings instance
settings = Settings()

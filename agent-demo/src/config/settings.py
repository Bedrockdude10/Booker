"""Application settings from environment variables."""

from pydantic_settings import BaseSettings
from pydantic import Field


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

    class Config:
        env_file = ".env"
        case_sensitive = False


# Global settings instance
settings = Settings()

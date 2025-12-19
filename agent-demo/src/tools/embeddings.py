"""Embedding generation using Sentence Transformers.

IMPORTANT: This module is for READ-ONLY operations (generating query embeddings).
For adding embeddings to database documents, use update_embeddings.py script.
"""

from sentence_transformers import SentenceTransformer

_EMBEDDING_MODEL: SentenceTransformer | None = None


def get_embedding_model() -> SentenceTransformer:
    """Lazy load embedding model (all-mpnet-base-v2, 768 dims)."""
    global _EMBEDDING_MODEL
    if _EMBEDDING_MODEL is None:
        _EMBEDDING_MODEL = SentenceTransformer('all-mpnet-base-v2')
    return _EMBEDDING_MODEL


def generate_embedding(text: str) -> list[float]:
    """
    Generate 768-dimensional embedding for text.

    Used for query-time embedding generation (user search queries).
    NOT for updating database documents - use update_embeddings.py for that.
    """
    model = get_embedding_model()
    embedding = model.encode([text], convert_to_numpy=True)
    return embedding[0].tolist()

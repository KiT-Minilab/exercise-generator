CREATE TYPE provider_enum AS ENUM ('OPENAI', 'GOOGLE');

CREATE TABLE IF NOT EXISTS exercise_generation_config (
    id SERIAL PRIMARY KEY,
    provider provider_enum NOT NULL,
    api_endpoint TEXT NOT NULL,
    api_key VARCHAR(128) NOT NULL,
    top_p REAL NOT NULL DEFAULT 1,
    temperature REAL NOT NULL DEFAULT 1,
    max_tokens INTEGER NOT NULL DEFAULT 256,
    messages JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
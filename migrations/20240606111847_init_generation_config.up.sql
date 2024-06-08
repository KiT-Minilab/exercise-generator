CREATE TYPE provider_enum AS ENUM ('OPENAI', 'GOOGLE');

CREATE TABLE IF NOT EXISTS generation_config (
    id SERIAL PRIMARY KEY,
    provider provider_enum UNIQUE NOT NULL,
    gen_model TEXT NOT NULL,
    top_p REAL NOT NULL DEFAULT 1,
    temperature REAL NOT NULL DEFAULT 1,
    max_tokens INTEGER NOT NULL DEFAULT 256,
    system_message TEXT NOT NULL DEFAULT '',
    assistant_message TEXT NOT NULL DEFAULT '',
    user_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
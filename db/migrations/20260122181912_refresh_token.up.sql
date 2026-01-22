PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    user_id INTEGER NOT NULL,

    token_hash TEXT NOT NULL UNIQUE,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    last_used_at DATETIME,

    expires_at DATETIME NOT NULL,

    revoked_at DATETIME,

    replaced_by_token_id INTEGER,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (replaced_by_token_id) REFERENCES refresh_tokens(id)
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id
    ON refresh_tokens(user_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at
    ON refresh_tokens(expires_at);


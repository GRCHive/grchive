CREATE TABLE user_sessions (
    session_id VARCHAR(36) NOT NULL UNIQUE,
    last_active_time TIMESTAMPTZ NOT NULL,
    expiration_time TIMESTAMPTZ NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    access_token TEXT NOT NULL,
    id_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

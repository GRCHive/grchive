CREATE TABLE user_sessions (
    session_id  VARCHAR(36) NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL,
    last_active_time TIMESTAMP NOT NULL,
    expiration_time TIMESTAMP NOT NULL,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    access_token TEXT NOT NULL,
    id_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL
);

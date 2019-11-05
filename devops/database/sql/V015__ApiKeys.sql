CREATE TABLE api_keys (
    hashed_api_key VARCHAR(128) NOT NULL UNIQUE,
    expiration_date TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE
);

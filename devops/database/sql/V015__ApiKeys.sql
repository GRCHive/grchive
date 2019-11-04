CREATE TABLE api_keys (
    hashed_api_key bytea NOT NULL,
    salt VARCHAR(256) NOT NULL,
    expiration_date TIMESTAMPTZ NOT NULL,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE
);

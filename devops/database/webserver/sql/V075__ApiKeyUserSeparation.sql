ALTER TABLE api_keys
ADD COLUMN id BIGSERIAL;

ALTER TABLE api_keys
ADD PRIMARY KEY (id);

CREATE TABLE api_key_to_users (
    api_key_id BIGINT NOT NULL UNIQUE REFERENCES api_keys(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(api_key_id, user_id)
);

INSERT INTO api_key_to_users (api_key_id, user_id)
SELECT key.id, key.user_id
FROM api_keys AS key;

ALTER TABLE api_keys
DROP COLUMN user_id;

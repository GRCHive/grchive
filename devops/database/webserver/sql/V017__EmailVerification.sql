CREATE TABLE email_verification (
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    code VARCHAR(36) NOT NULL,
    verification_sent TIMESTAMPTZ,
    verification_received TIMESTAMPTZ
);

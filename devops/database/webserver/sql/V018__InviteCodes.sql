CREATE TABLE invitation_codes (
    id BIGSERIAL,
    from_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    to_email VARCHAR(320) NOT NULL,
    sent_time TIMESTAMPTZ,
    used_time TIMESTAMPTZ,
    PRIMARY KEY(id, from_org_id),
    FOREIGN KEY(from_user_id, from_org_id)
        REFERENCES user_orgs(user_id, org_id)
        ON DELETE CASCADE
);

-- Prevent an end user being spammed with invitation emails by forcing only a single
-- pending invitation to a user from an organization at any given moment.
CREATE UNIQUE INDEX invitation_codes_null_used_idx ON invitation_codes (from_org_id, to_email)
WHERE used_time IS NULL;

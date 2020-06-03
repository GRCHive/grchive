CREATE TABLE db_refresh_diff_message_recipients (
    db_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    FOREIGN KEY(db_id, org_id) REFERENCES database_resources(id, org_id) ON DELETE CASCADE
);

CREATE INDEX ON db_refresh_diff_message_recipients(db_id);

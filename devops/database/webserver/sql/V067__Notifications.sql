CREATE TABLE notifications (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    time TIMESTAMPTZ NOT NULL,
    subject TEXT NOT NULL,
    subject_data JSONB NOT NULL,
    verb TEXT NOT NULL,
    object TEXT NOT NULL,
    object_data JSONB NOT NULL,
    PRIMARY KEY(org_id, id)
);
CREATE INDEX ON notifications(org_id, id);

CREATE TABLE user_notifications (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    notification_id BIGINT NOT NULL,
    read TIMESTAMPTZ,
    sent TIMESTAMPTZ,
    FOREIGN KEY(org_id, notification_id) REFERENCES notifications(org_id, id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX ON user_notifications(org_id, user_id, notification_id);
CREATE INDEX ON user_notifications(org_id, user_id);

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    time TIMESTAMPTZ NOT NULL,
    subject TEXT NOT NULL,
    subject_data JSONB NOT NULL,
    verb TEXT NOT NULL,
    object TEXT NOT NULL,
    object_data JSONB NOT NULL
);
CREATE INDEX ON notifications(org_id);

CREATE TABLE user_notifications (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    notification_id BIGINT NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    read TIMESTAMPTZ,
    sent TIMESTAMPTZ
);
CREATE UNIQUE INDEX ON user_notifications(org_id, user_id, notification_id);
CREATE INDEX ON user_notifications(org_id, user_id);

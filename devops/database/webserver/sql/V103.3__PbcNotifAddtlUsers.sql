CREATE TABLE org_pbc_notification_cadence_addtl_users (
    cadence_id BIGINT NOT NULL REFERENCES org_pbc_notification_cadence_settings(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ON org_pbc_notification_cadence_addtl_users(cadence_id);

CREATE TABLE integration_system_link (
    integration_id BIGINT NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    system_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(system_id, org_id) REFERENCES systems(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON integration_system_link(system_id);

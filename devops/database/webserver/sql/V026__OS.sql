CREATE TABLE operating_systems (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    os_name VARCHAR(256) NOT NULL,
    os_version VARCHAR(256) NOT NULL,
    PRIMARY KEY(id, org_id)
);

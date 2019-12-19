CREATE TABLE servers (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    os_id BIGINT,
    external_ip VARCHAR(40),
    host VARCHAR(256),
    location VARCHAR(256),
    description TEXT,
    FOREIGN KEY(os_id, org_id) REFERENCES operating_systems(id, org_id) ON DELETE NO ACTION
);

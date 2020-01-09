CREATE TABLE infrastructure_servers (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    ip_address TEXT,
    operating_system TEXT,
    location TEXT,
    PRIMARY KEY (id, org_id),
    UNIQUE(org_id, name)
);

CREATE TABLE systems (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    purpose TEXT,
    description TEXT,
    PRIMARY KEY(id, org_id)
);

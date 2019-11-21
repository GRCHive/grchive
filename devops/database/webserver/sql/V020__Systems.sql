CREATE TABLE systems (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id),
    name TEXT NOT NULL,
    purpose TEXT,
    description TEXT,
    PRIMARY KEY(id, org_id)
);

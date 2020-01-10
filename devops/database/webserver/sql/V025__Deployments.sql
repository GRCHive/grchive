CREATE TABLE deployments (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    deployment_type INTEGER NOT NULL,
    PRIMARY KEY (id, org_id)
);

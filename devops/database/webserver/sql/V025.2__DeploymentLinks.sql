CREATE TABLE deployment_system_link (
    system_id BIGINT NOT NULL UNIQUE,
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (system_id, org_id)
        REFERENCES systems(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE
);

CREATE TABLE deployment_db_link (
    db_id BIGINT NOT NULL UNIQUE,
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (db_id, org_id)
        REFERENCES database_resources(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE
);

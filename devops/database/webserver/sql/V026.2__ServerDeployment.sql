CREATE TABLE deployment_server_link (
    server_id BIGINT NOT NULL,
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (server_id, org_id)
        REFERENCES infrastructure_servers(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE,
    UNIQUE(server_id, deployment_id)
);


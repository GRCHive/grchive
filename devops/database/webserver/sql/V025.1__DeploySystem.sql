CREATE TABLE deployment_system_link (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
) INHERITS (_base_deployment_link);

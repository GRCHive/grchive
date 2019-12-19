CREATE TABLE deployment_db_link (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE
) INHERITS (_base_deployment_link);

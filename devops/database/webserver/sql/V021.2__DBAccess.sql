CREATE TABLE resource_database_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

CREATE TABLE resource_db_conn_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

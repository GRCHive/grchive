CREATE TABLE resource_db_sql_query_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

CREATE UNIQUE INDEX ON resource_db_sql_access(role_id);

INSERT INTO resource_db_sql_query_access (role_id, org_id, access_type)
SELECT id, org_id, 7
FROM organization_available_roles
WHERE is_admin_role = true;

INSERT INTO resource_db_sql_query_access (role_id, org_id, access_type)
SELECT id, org_id, 1
FROM organization_available_roles
WHERE is_admin_role = false;

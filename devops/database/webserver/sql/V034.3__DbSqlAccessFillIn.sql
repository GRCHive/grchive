INSERT INTO resource_db_sql_access (role_id, org_id, access_type)
SELECT id, org_id, 7
FROM organization_available_roles
WHERE is_admin_role = true;

INSERT INTO resource_db_sql_access (role_id, org_id, access_type)
SELECT id, org_id, 1
FROM organization_available_roles
WHERE is_admin_role = false;

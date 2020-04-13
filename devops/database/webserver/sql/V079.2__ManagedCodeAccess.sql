CREATE TABLE resource_managed_code_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

CREATE UNIQUE INDEX ON resource_managed_code_access(role_id);

INSERT INTO resource_managed_code_access (role_id, org_id, access_type)
SELECT id, org_id, 7
FROM organization_available_roles
WHERE is_admin_role = true;

INSERT INTO resource_managed_code_access (role_id, org_id, access_type)
SELECT id, org_id, 1
FROM organization_available_roles
WHERE is_admin_role = false;

ALTER TABLE managed_kotlin_code
RENAME TO managed_code;
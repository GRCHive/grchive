CREATE TABLE resource_deployment_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

INSERT INTO resource_deployment_access (role_id, org_id, access_type)
SELECT id, org_id, 7
FROM organization_available_roles;

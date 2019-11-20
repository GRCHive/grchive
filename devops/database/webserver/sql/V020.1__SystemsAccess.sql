CREATE TABLE resource_systems_access () INHERITS (_base_resource_access);

INSERT INTO resource_systems_access (role_id, access_type)
SELECT r.id, 7
FROM organization_available_roles AS r;

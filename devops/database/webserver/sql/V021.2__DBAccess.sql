CREATE TABLE resource_database_access () INHERITS (_base_resource_access);

INSERT INTO resource_database_access (role_id, access_type)
SELECT r.id, 7
FROM organization_available_roles AS r;

CREATE TABLE resource_db_conn_access () INHERITS (_base_resource_access);

INSERT INTO resource_db_conn_access (role_id, access_type)
SELECT r.id, 7
FROM organization_available_roles AS r;

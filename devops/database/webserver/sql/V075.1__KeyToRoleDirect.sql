CREATE TABLE api_key_roles (
    role_id BIGINT NOT NULL,
    api_key_id BIGINT NOT NULL UNIQUE REFERENCES api_keys(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id, org_id) REFERENCES organization_available_roles(id, org_id) ON DELETE RESTRICT,
    UNIQUE(role_id, api_key_id, org_id)
);

CREATE INDEX ON api_key_roles (role_id, org_id);

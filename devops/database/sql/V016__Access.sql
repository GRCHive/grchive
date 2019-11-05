CREATE TABLE organization_available_roles (
    id BIGSERIAL PRIMARY KEY,
    is_default_role BOOLEAN NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    UNIQUE(name, org_id)
);

CREATE UNIQUE INDEX org_roles_default_index ON organization_available_roles (is_default_role)
WHERE is_default_role = true;

CREATE UNIQUE INDEX org_roles_id_org_index ON organization_available_roles (id, org_id);

CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL UNIQUE REFERENCES organization_available_roles(id) ON DELETE CASCADE,
    permissions JSONB
);

CREATE TABLE user_roles (
    role_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id INTEGER NOT NULL,
    UNIQUE(user_id, org_id),
    CONSTRAINT role_org_key
        FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE,
    CONSTRAINT user_org_key
        FOREIGN KEY(user_id, org_id)
        REFERENCES users(id, org_id)
        ON DELETE CASCADE
);

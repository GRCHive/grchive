CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    -- group id on okta
    org_group_id VARCHAR(256) UNIQUE,
    -- group name on okta
    org_group_name VARCHAR(256) UNIQUE,
    -- our stored for the organization
    org_name TEXT NOT NULL UNIQUE,
    saml_iden VARCHAR(100) UNIQUE REFERENCES saml_idp(idpIdenOkta) ON DELETE NO ACTION
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    email VARCHAR(320) NOT NULL UNIQUE
);

CREATE TABLE user_orgs (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    org_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, org_id)
);

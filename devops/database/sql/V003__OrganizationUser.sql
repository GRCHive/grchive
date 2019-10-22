CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    -- group id on okta
    org_group_id VARCHAR(256) NOT NULL UNIQUE,
    -- group name on okta
    org_group_name VARCHAR(256) NOT NULL UNIQUE,
    -- our stored for the organization
    org_name TEXT NOT NULL UNIQUE,
    saml_iden VARCHAR(100) NOT NULL UNIQUE REFERENCES saml_idp(idpIdenOkta) ON DELETE NO ACTION
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    okta_id VARCHAR(256) NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    email VARCHAR(320) NOT NULL UNIQUE
);

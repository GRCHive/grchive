CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    -- group id on okta
    org_group_id VARCHAR(256) NOT NULL UNIQUE,
    -- group name on okta
    org_group_name VARCHAR(256) NOT NULL,
    -- our stored for the organization
    org_name TEXT NOT NULL
);

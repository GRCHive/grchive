CREATE TABLE deployment_options (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    description TEXT
);

CREATE TABLE deployments (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    option_id INTEGER REFERENCES deployment_options(id) ON DELETE RESTRICT,
    vendor_solution VARCHAR(256),
    description TEXT,
    PRIMARY KEY(id, org_id)
);

CREATE TABLE _base_deployment_link (
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL
);

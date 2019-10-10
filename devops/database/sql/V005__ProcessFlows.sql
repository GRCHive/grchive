CREATE TABLE process_flows (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    description TEXT,
    created_time TIMESTAMPTZ NOT NULL,
    last_updated_time TIMESTAMPTZ NOT NULL
);

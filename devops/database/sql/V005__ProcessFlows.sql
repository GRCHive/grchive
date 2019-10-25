CREATE TABLE process_flows (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    description TEXT,
    created_time TIMESTAMPTZ NOT NULL,
    last_updated_time TIMESTAMPTZ NOT NULL,
    UNIQUE(name, org_id)
);

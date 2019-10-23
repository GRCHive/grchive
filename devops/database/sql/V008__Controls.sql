CREATE TABLE process_flow_control_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE process_flow_controls (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    control_type INTEGER NOT NULL REFERENCES process_flow_control_types(id) ON DELETE RESTRICT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    frequency_type INTEGER NOT NULL,
    frequency_interval INTEGER NOT NULL,
    owner_id BIGINT REFERENCES users(id) ON DELETE NO ACTION
);

CREATE TABLE process_flow_control_node (
    control_id BIGINT NOT NULL REFERENCES process_flow_controls(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE
);

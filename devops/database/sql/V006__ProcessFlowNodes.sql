CREATE TABLE process_flow_node_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE process_flow_nodes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256),
    process_flow_id INTEGER NOT NULL REFERENCES process_flows(id) ON DELETE RESTRICT,
    description TEXT,
    node_type INTEGER NOT NULL REFERENCES process_flow_node_types(id) ON DELETE RESTRICT
);

CREATE TABLE process_flow_node_inputs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    parent_node_id INTEGER NOT NULL REFERENCES process_flow_nodes(id)
);

CREATE TABLE process_flow_node_outputs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    parent_node_id INTEGER NOT NULL REFERENCES process_flow_nodes(id)
);

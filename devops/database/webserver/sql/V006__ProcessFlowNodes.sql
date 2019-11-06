CREATE TABLE process_flow_node_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE process_flow_nodes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256),
    process_flow_id BIGINT NOT NULL REFERENCES process_flows(id) ON DELETE CASCADE,
    description TEXT,
    node_type INTEGER NOT NULL REFERENCES process_flow_node_types(id) ON DELETE RESTRICT
);

CREATE TABLE process_flow_input_output_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE TABLE process_flow_node_inputs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    parent_node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    io_type_id INTEGER NOT NULL REFERENCES process_flow_input_output_type(id) ON DELETE CASCADE
);

CREATE TABLE process_flow_node_outputs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    parent_node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    io_type_id INTEGER NOT NULL REFERENCES process_flow_input_output_type(id) ON DELETE CASCADE
);

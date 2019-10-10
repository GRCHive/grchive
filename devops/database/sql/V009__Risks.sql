CREATE TABLE process_flow_risks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE RESTRICT,
    control_id BIGINT REFERENCES process_flow_controls(id) ON DELETE RESTRICT
);

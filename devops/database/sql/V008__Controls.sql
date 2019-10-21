CREATE TABLE process_flow_controls (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE TABLE process_flow_control_node (
    control_id BIGINT NOT NULL REFERENCES process_flow_controls(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE
);

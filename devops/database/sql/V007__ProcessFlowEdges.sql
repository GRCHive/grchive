CREATE TABLE process_flow_edges (
    id BIGSERIAL PRIMARY KEY,
    input_id BIGINT NOT NULL REFERENCES process_flow_node_inputs(id) ON DELETE CASCADE,
    output_id BIGINT NOT NULL REFERENCES process_flow_node_outputs(id) ON DELETE CASCADE
);

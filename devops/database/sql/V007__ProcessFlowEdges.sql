CREATE TABLE process_flow_edges (
    id BIGSERIAL PRIMARY KEY,
    process_flow_id BIGINT NOT NULL REFERENCES process_flows(id) ON DELETE RESTRICT,
    input_id BIGINT NOT NULL REFERENCES process_flow_node_inputs(id) ON DELETE RESTRICT,
    output_id BIGINT NOT NULL REFERENCES process_flow_node_outputs(id) ON DELETE RESTRICT
);

CREATE TABLE process_flow_edges (
    id SERIAL PRIMARY KEY,
    process_flow_id INTEGER NOT NULL REFERENCES process_flows(id) ON DELETE RESTRICT,
    input_id INTEGER NOT NULL REFERENCES process_flow_node_inputs(id) ON DELETE RESTRICT,
    output_id INTEGER NOT NULL REFERENCES process_flow_node_outputs(id) ON DELETE RESTRICT
);

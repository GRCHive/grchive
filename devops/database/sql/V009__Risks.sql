CREATE TABLE process_flow_risks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    node_id INTEGER NOT NULL REFERENCES process_flow_nodes(id) ON DELETE RESTRICT,
    control_id INTEGER REFERENCES process_flow_controls(id) ON DELETE RESTRICT
);

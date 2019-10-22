CREATE TABLE process_flow_risks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT
);

CREATE TABLE process_flow_risk_node (
    risk_id BIGINT NOT NULL REFERENCES process_flow_risks(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    UNIQUE(risk_id, node_id)
);

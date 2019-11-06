CREATE TABLE process_flow_risks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    UNIQUE(name, org_id)
);

CREATE TABLE process_flow_risk_node (
    risk_id BIGINT NOT NULL REFERENCES process_flow_risks(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    UNIQUE(risk_id, node_id)
);

CREATE TABLE process_flow_risk_control (
    risk_id BIGINT NOT NULL REFERENCES process_flow_risks(id) ON DELETE CASCADE,
    control_id BIGINT NOT NULL REFERENCES process_flow_controls(id) ON DELETE CASCADE,
    UNIQUE(risk_id, control_id)
);

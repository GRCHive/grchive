CREATE TABLE node_system_link (
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    system_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(system_id, org_id) REFERENCES systems(id, org_id) ON DELETE CASCADE
);

CREATE INDEX ON node_system_link(org_id, system_id, node_id);

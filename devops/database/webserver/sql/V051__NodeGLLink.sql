CREATE TABLE node_gl_link (
    node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE,
    gl_account_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(gl_account_id, org_id) REFERENCES general_ledger_accounts(id, org_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON node_gl_link(org_id, gl_account_id, node_id);

ALTER TABLE process_flow_node_types
ADD COLUMN can_link_to_gl BOOLEAN NOT NULL DEFAULT false;

UPDATE process_flow_node_types
SET can_link_to_gl = true
WHERE name = 'General Ledger';

ALTER TABLE process_flow_node_types
ALTER COLUMN can_link_to_gl DROP DEFAULT;


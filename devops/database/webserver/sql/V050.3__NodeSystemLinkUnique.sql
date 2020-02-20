DROP INDEX node_system_link_org_id_system_id_node_id_idx;
CREATE UNIQUE INDEX ON node_system_link(org_id, system_id, node_id);

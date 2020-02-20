ALTER TABLE process_flow_node_types
ADD COLUMN can_link_to_system BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE process_flow_node_types
ALTER COLUMN can_link_to_system DROP DEFAULT;

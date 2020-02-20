UPDATE process_flow_node_types
SET can_link_to_system = true
WHERE name = 'System';

CREATE TABLE process_flow_node_display_settings (
    id SERIAL PRIMARY KEY,
    parent_node_id BIGINT NOT NULL REFERENCES process_flow_nodes(id) ON DELETE CASCADE UNIQUE,
    settings JSONB
);

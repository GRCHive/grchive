CREATE OR REPLACE FUNCTION audit_process_flow_node_outputs_change(r process_flow_node_outputs, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        org_id INTEGER;
        flow_id BIGINT;
        event_id BIGINT;
    BEGIN
        SELECT n.process_flow_id INTO flow_id
        FROM process_flow_nodes AS n
        INNER JOIN process_flow_node_outputs AS io
            ON io.parent_node_id = n.id
        WHERE n.id = r.id;
        
        SELECT p.org_id INTO org_id
        FROM process_flows AS p
        WHERE p.id = flow_id;

        SELECT generic_audit_event(org_id,
            'process_flow_node_outputs',
            r.id,
            jsonb_build_object(
                'process_flow_id', flow_id,
                'node_id', r.parent_node_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION audit_process_flow_node_inputs_change(r process_flow_node_inputs, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        org_id INTEGER;
        flow_id BIGINT;
        event_id BIGINT;
    BEGIN
        SELECT n.process_flow_id INTO flow_id
        FROM process_flow_nodes AS n
        INNER JOIN process_flow_node_inputs AS io
            ON io.parent_node_id = n.id
        WHERE io.id = r.id;

        IF flow_id IS NULL THEN
            -- This means node is already deleted. We need to rely on
            -- using the audit log instead. We probably can't rely on
            -- the process flow existing either.
            SELECT CAST(aud.resource_extra_data #>> '{process_flow_id}' AS BIGINT), aud.org_id INTO flow_id, org_id
            FROM global_audit_event_history AS aud
            WHERE aud.resource_type = 'process_flow_nodes' AND CAST(aud.resource_id AS BIGINT) = r.parent_node_id;
        ELSE
            SELECT p.org_id INTO org_id
            FROM process_flows AS p
            WHERE p.id = flow_id;
        END IF;

        SELECT generic_audit_event(org_id,
            'process_flow_node_inputs',
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
        WHERE io.id = r.id;
        
        IF flow_id IS NULL THEN
            -- This means node is already deleted. We need to rely on
            -- using the audit log instead. We probably can't rely on
            -- the process flow existing either.
            SELECT CAST(aud.resource_extra_data #>> '{process_flow_id}' AS BIGINT), aud.org_id INTO flow_id, org_id
            FROM global_audit_event_history AS aud
            WHERE aud.resource_type = 'process_flow_nodes' AND CAST(aud.resource_id AS BIGINT) = r.parent_node_id;
        ELSE
            SELECT p.org_id INTO org_id
            FROM process_flows AS p
            WHERE p.id = flow_id;
        END IF;

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

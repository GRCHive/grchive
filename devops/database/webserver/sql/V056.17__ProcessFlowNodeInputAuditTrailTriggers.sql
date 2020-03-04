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

        SELECT p.org_id INTO org_id
        FROM process_flows AS p
        WHERE p.id = flow_id;

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

CREATE OR REPLACE FUNCTION insert_audit_process_flow_node_inputs_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_node_inputs_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_process_flow_node_inputs_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_node_inputs_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_process_flow_node_inputs_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_node_inputs_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_process_flow_node_inputs
    AFTER INSERT ON process_flow_node_inputs
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_process_flow_node_inputs_change();

CREATE TRIGGER trigger_update_process_flow_node_inputs
    AFTER UPDATE ON process_flow_node_inputs
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_process_flow_node_inputs_change();

CREATE TRIGGER trigger_delete_process_flow_node_inputs
    BEFORE DELETE ON process_flow_node_inputs
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_process_flow_node_inputs_change();

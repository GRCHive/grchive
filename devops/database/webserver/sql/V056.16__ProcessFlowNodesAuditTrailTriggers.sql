CREATE OR REPLACE FUNCTION audit_process_flow_nodes_change(r process_flow_nodes, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        org_id INTEGER;
        event_id BIGINT;
    BEGIN
        SELECT p.org_id INTO org_id
        FROM process_flows AS p
        INNER JOIN process_flow_nodes AS n
            ON n.process_flow_id = p.id
        WHERE n.id = r.id;

        SELECT generic_audit_event(org_id,
            'process_flow_nodes',
            r.id,
            jsonb_build_object('process_flow_id', r.process_flow_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_process_flow_nodes_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_nodes_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_process_flow_nodes_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_nodes_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_process_flow_nodes_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_nodes_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_process_flow_nodes
    AFTER INSERT ON process_flow_nodes
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_process_flow_nodes_change();

CREATE TRIGGER trigger_update_process_flow_nodes
    AFTER UPDATE ON process_flow_nodes
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_process_flow_nodes_change();

CREATE TRIGGER trigger_delete_process_flow_nodes
    BEFORE DELETE ON process_flow_nodes
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_process_flow_nodes_change();

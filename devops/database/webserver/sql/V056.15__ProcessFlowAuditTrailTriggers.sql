CREATE OR REPLACE FUNCTION audit_process_flows_change(r process_flows, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'process_flows',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_process_flows_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flows_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_process_flows_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flows_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_process_flows_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flows_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_process_flows
    AFTER INSERT ON process_flows
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_process_flows_change();

CREATE TRIGGER trigger_update_process_flows
    AFTER UPDATE ON process_flows
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_process_flows_change();

CREATE TRIGGER trigger_delete_process_flows
    BEFORE DELETE ON process_flows
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_process_flows_change();

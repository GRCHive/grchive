CREATE OR REPLACE FUNCTION audit_process_flow_risks_change(r process_flow_risks, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'process_flow_risks',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_process_flow_risks_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_risks_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_process_flow_risks_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_risks_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_process_flow_risks_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_process_flow_risks_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_process_flow_risks
    AFTER INSERT ON process_flow_risks
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_process_flow_risks_change();

CREATE TRIGGER trigger_update_process_flow_risks
    AFTER UPDATE ON process_flow_risks
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_process_flow_risks_change();

CREATE TRIGGER trigger_delete_process_flow_risks
    BEFORE DELETE ON process_flow_risks
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_process_flow_risks_change();
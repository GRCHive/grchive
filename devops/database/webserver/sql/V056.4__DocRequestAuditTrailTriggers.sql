CREATE OR REPLACE FUNCTION audit_document_request_change(r document_requests, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'document_requests',
            r.id,
            '{}'::jsonb,
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_document_request_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_document_request_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_document_request_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_document_request_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_document_request_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_document_request_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_document_requests
    AFTER INSERT ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_document_request_change();

CREATE TRIGGER trigger_update_document_requests
    AFTER UPDATE ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_document_request_change();

CREATE TRIGGER trigger_delete_document_requests
    BEFORE DELETE ON document_requests
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_document_request_change();

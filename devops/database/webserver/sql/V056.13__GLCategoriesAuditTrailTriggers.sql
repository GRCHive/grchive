CREATE OR REPLACE FUNCTION audit_general_ledger_categories_change(r general_ledger_categories, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'general_ledger_categories',
            r.id,
            jsonb_build_object('gl_parent_cat_id', r.parent_category_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_general_ledger_categories_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_general_ledger_categories_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_general_ledger_categories_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_general_ledger_categories_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_general_ledger_categories_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_general_ledger_categories_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_general_ledger_categories
    AFTER INSERT ON general_ledger_categories
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_general_ledger_categories_change();

CREATE TRIGGER trigger_update_general_ledger_categories
    AFTER UPDATE ON general_ledger_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_general_ledger_categories_change();

CREATE TRIGGER trigger_delete_general_ledger_categories
    BEFORE DELETE ON general_ledger_categories
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_general_ledger_categories_change();

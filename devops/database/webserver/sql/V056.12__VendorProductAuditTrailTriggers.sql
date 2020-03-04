CREATE OR REPLACE FUNCTION audit_vendor_products_change(r vendor_products, action VARCHAR(64))
    RETURNS void AS
$$
    DECLARE
        event_id BIGINT;
    BEGIN
        SELECT generic_audit_event(r.org_id,
            'vendor_products',
            r.id,
            jsonb_build_object('vendor_id', r.vendor_id),
            action
        ) INTO event_id;
        INSERT INTO audit_resource_modifications(event_id, data)
        VALUES (event_id, to_jsonb(r));
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_audit_vendor_products_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendor_products_change(NEW, 'INSERT');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_audit_vendor_products_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendor_products_change(NEW, 'UPDATE');
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_audit_vendor_products_change()
    RETURNS trigger AS
$$
    BEGIN
        PERFORM audit_vendor_products_change(OLD, 'DELETE');
        RETURN OLD;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_insert_vendor_products
    AFTER INSERT ON vendor_products
    FOR EACH ROW
    EXECUTE FUNCTION insert_audit_vendor_products_change();

CREATE TRIGGER trigger_update_vendor_products
    AFTER UPDATE ON vendor_products
    FOR EACH ROW
    EXECUTE FUNCTION update_audit_vendor_products_change();

CREATE TRIGGER trigger_delete_vendor_products
    BEFORE DELETE ON vendor_products
    FOR EACH ROW
    EXECUTE FUNCTION delete_audit_vendor_products_change();

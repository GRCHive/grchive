CREATE OR REPLACE FUNCTION update_general_ledger_categories_prevent_loop()
    RETURNS trigger AS
$$
    DECLARE
        is_loop BOOLEAN;
    BEGIN
        WITH RECURSIVE child_categories AS (
            SELECT id
            FROM general_ledger_categories
            WHERE parent_category_id = OLD.id
            UNION
                SELECT c.id
                FROM general_ledger_categories AS c
                INNER JOIN child_categories AS ch ON ch.id = c.parent_category_id
        )
        SELECT NEW.parent_category_id = ANY(ARRAY_AGG(id)) INTO is_loop
        FROM child_categories;

        IF is_loop THEN
            RAISE EXCEPTION 'Loop in GL categories.';
        END IF;

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_general_ledger_categories_prevent_loop ON general_ledger_categories;
CREATE TRIGGER trigger_update_general_ledger_categories_prevent_loop
    BEFORE UPDATE ON general_ledger_categories
    FOR EACH ROW
    EXECUTE FUNCTION update_general_ledger_categories_prevent_loop();

CREATE OR REPLACE FUNCTION create_default_control_folder()
    RETURNS trigger AS
$$
    DECLARE
        folder_id BIGINT;
    BEGIN
        INSERT INTO file_folders (org_id, name)
        VALUES (NEW.org_id, 'Default')
        RETURNING id INTO folder_id;

        INSERT INTO control_folder_link (control_id, folder_id, org_id)
        VALUES (NEW.id, folder_id, NEW.org_id);

        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_default_control_folder ON process_flow_controls;
CREATE TRIGGER trigger_default_control_folder
    AFTER INSERT ON process_flow_controls
    FOR EACH ROW
    EXECUTE FUNCTION create_default_control_folder();


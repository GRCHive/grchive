CREATE TABLE file_folders (
    id BIGSERIAL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name TEXT,
    PRIMARY KEY(id, org_id)
);
CREATE INDEX ON file_folders(org_id);

CREATE TABLE control_folder_link (
    control_id BIGINT NOT NULL,
    folder_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(control_id, org_id) REFERENCES process_flow_controls(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(folder_id, org_id) REFERENCES file_folders(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON control_folder_link(org_id, folder_id, control_id);

CREATE TABLE file_folder_link (
    folder_id BIGINT NOT NULL,
    file_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(file_id, org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(folder_id, org_id) REFERENCES file_folders(id, org_id) ON DELETE CASCADE
);
CREATE INDEX ON file_folder_link(org_id, folder_id, file_id);

CREATE OR REPLACE FUNCTION migrate_control_doc_cats_to_folders(input_control_id BIGINT, input_org_id INTEGER)
    RETURNS void AS
$$
    DECLARE
        input_folder_id INTEGER;
        output_folder_id INTEGER;
    BEGIN
        INSERT INTO file_folders (org_id, name)
        SELECT input_org_id, 'Input'
        FROM process_flow_controls
        WHERE id = input_control_id
        RETURNING id INTO input_folder_id;

        INSERT INTO file_folders (org_id, name)
        SELECT input_org_id, 'Output'
        FROM process_flow_controls
        WHERE id = input_control_id
        RETURNING id INTO output_folder_id;

        INSERT INTO control_folder_link (control_id, folder_id, org_id)
        VALUES 
            (input_control_id, input_folder_id, input_org_id),
            (input_control_id, output_folder_id, input_org_id);

        INSERT INTO file_folder_link (folder_id, file_id, org_id)
        SELECT input_folder_id, file.id, input_org_id
        FROM file_metadata AS file
        INNER JOIN process_flow_control_documentation_categories AS cat
            ON cat.id = file.category_id
        INNER JOIN controls_input_documentation AS inp
            ON inp.category_id = cat.id
        WHERE inp.control_id = input_control_id AND inp.org_id = input_org_id;

        INSERT INTO file_folder_link (folder_id, file_id, org_id)
        SELECT output_folder_id, file.id, input_org_id
        FROM file_metadata AS file
        INNER JOIN process_flow_control_documentation_categories AS cat
            ON cat.id = file.category_id
        INNER JOIN controls_output_documentation AS ot
            ON ot.category_id = cat.id
        WHERE ot.control_id = input_control_id AND ot.org_id = input_org_id;
    END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    PERFORM migrate_control_doc_cats_to_folders(id, org_id)
    FROM process_flow_controls;
END $$ LANGUAGE plpgsql;

DROP TABLE controls_input_documentation;
DROP TABLE controls_output_documentation;
DROP TABLE _base_link_controls_control_documentation;

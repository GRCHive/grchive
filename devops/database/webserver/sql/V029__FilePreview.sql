CREATE TABLE file_previews (
    original_file_id BIGINT NOT NULL UNIQUE,
    preview_file_id BIGINT UNIQUE,
    category_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(original_file_id, category_id, org_id)
        REFERENCES process_flow_control_documentation_file(id, category_id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY(preview_file_id, category_id, org_id)
        REFERENCES process_flow_control_documentation_file(id, category_id, org_id)
        ON DELETE CASCADE
);

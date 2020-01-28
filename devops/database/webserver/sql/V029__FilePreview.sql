CREATE TABLE file_previews (
    file_id BIGINT NOT NULL,
    original_storage_id BIGINT NOT NULL UNIQUE,
    preview_storage_id BIGINT UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(file_id, org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(original_storage_id, org_id)
        REFERENCES file_storage(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY(preview_storage_id, org_id)
        REFERENCES file_storage(id, org_id)
        ON DELETE CASCADE
);

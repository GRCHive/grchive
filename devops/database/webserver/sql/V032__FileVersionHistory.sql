CREATE TABLE file_version_history (
    file_id BIGINT NOT NULL,
    file_storage_id BIGINT NOT NULL UNIQUE,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    FOREIGN KEY(file_id, org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(file_storage_id, org_id) REFERENCES file_storage(id, org_id) ON DELETE CASCADE
);

CREATE INDEX file_version_history_file_id_index ON file_version_history(file_id);

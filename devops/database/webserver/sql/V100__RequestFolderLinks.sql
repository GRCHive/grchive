CREATE TABLE request_folder_link (
    request_id BIGINT NOT NULL,
    folder_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(folder_id, org_id) REFERENCES file_folders(id, org_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON request_folder_link(org_id, request_id);
CREATE INDEX ON request_folder_link(folder_id);

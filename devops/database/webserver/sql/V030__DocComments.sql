CREATE TABLE file_comments (
    file_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY(file_id, org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE
)

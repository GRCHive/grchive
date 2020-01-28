CREATE TABLE document_request_comments (
    request_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE
);

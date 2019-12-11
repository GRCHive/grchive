CREATE TABLE document_request_comments (
    request_id BIGINT NOT NULL,
    cat_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, cat_id, org_id) REFERENCES document_requests(id, cat_id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE
);

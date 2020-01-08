CREATE TABLE deployment_request_link (
    deployment_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    request_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    FOREIGN KEY (deployment_id, org_id)
        REFERENCES deployments(id, org_id)
        ON DELETE CASCADE,
    FOREIGN KEY (request_id, category_id, org_id)
        REFERENCES document_requests(id, cat_id, org_id)
        ON DELETE CASCADE,
    UNIQUE(deployment_id, request_id)
);

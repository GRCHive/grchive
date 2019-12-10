CREATE TABLE document_requests (
    id BIGSERIAL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    cat_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    requested_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id, cat_id, org_id),
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE
);

CREATE TABLE document_request_fulfillment (
    id BIGSERIAL,
    cat_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    fulfilled_file_id BIGINT,
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(fulfilled_file_id, cat_id, org_id) REFERENCES process_flow_control_documentation_file(id, category_id, org_id) ON DELETE NO ACTION
);

CREATE TABLE document_request_comments (
    id BIGSERIAL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT,
    post_time TIMESTAMPTZ NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    cat_id BIGINT NOT NULL,
    request_id BIGINT NOT NULL,
    PRIMARY KEY(id, request_id, org_id),
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, cat_id, org_id) REFERENCES document_requests(id, cat_id, org_id) ON DELETE CASCADE
);

CREATE TABLE resource_doc_request_access () INHERITS (_base_resource_access);
CREATE TABLE document_requests (
    id BIGSERIAL,
    name VARCHAR(256) NOT NULL,
    description TEXT,
    cat_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    requested_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE
);

CREATE INDEX document_requests_cat_index ON document_requests(cat_id);

CREATE TABLE document_request_fulfillment (
    id BIGSERIAL,
    request_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    fulfilled_file_id BIGINT,
    PRIMARY KEY(id, org_id),
    FOREIGN KEY(fulfilled_file_id,  org_id) REFERENCES file_metadata(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE
);

CREATE TABLE resource_doc_request_access (
    FOREIGN KEY(org_id)
        REFERENCES organizations(id)
        ON DELETE CASCADE,
    FOREIGN KEY(role_id, org_id)
        REFERENCES organization_available_roles(id, org_id)
        ON DELETE CASCADE
) INHERITS (_base_resource_access);

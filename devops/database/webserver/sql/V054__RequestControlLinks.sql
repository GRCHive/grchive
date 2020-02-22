CREATE TABLE request_control_link (
    request_id BIGINT NOT NULL,
    control_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(control_id, org_id) REFERENCES process_flow_controls(id, org_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON request_control_link(org_id, request_id);
CREATE UNIQUE INDEX ON request_control_link(org_id, control_id, request_id);

CREATE UNIQUE INDEX ON request_doc_cat_link(org_id, request_id);

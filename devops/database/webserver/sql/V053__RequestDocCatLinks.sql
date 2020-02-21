CREATE TABLE request_doc_cat_link (
    request_id BIGINT NOT NULL,
    cat_id BIGINT NOT NULL,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id, org_id) REFERENCES document_requests(id, org_id) ON DELETE CASCADE,
    FOREIGN KEY(cat_id, org_id) REFERENCES process_flow_control_documentation_categories(id, org_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX ON request_doc_cat_link(org_id, cat_id, request_id);

INSERT INTO request_doc_cat_link (request_id, cat_id, org_id)
SELECT req.id, req.cat_id, req.org_id
FROM document_requests AS req;

DROP INDEX document_requests_cat_index;
ALTER TABLE document_requests
DROP COLUMN cat_id;

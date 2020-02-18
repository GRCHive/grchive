DROP INDEX file_comments_file_id_org_id_idx;
DROP INDEX document_request_comments_request_id_org_id_idx;

CREATE INDEX ON document_request_comments(request_id, org_id);
CREATE INDEX ON file_comments(file_id, org_id);

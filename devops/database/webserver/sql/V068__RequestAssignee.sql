ALTER TABLE document_requests
ADD COLUMN assignee BIGINT;

ALTER TABLE document_requests
ADD CONSTRAINT document_requests_assignee_fkey FOREIGN KEY (assignee) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE database_sql_query_requests
ADD COLUMN assignee BIGINT;

ALTER TABLE database_sql_query_requests
ADD CONSTRAINT database_sql_query_requests_requests_assignee_fkey FOREIGN KEY (assignee) REFERENCES users(id) ON DELETE CASCADE;

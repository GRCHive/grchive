ALTER TABLE document_requests
ADD COLUMN due_date TIMESTAMPTZ;

ALTER TABLE database_sql_query_requests
ADD COLUMN due_date TIMESTAMPTZ;

DROP TABLE document_request_comments CASCADE;
DROP TABLE file_comments CASCADE;
DROP TABLE sql_request_comments CASCADE;

ALTER TABLE comments
ALTER COLUMN thread_id SET NOT NULL;

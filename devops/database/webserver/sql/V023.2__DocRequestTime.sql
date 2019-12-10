ALTER TABLE document_requests
ADD COLUMN request_time TIMESTAMPTZ;

ALTER TABLE document_requests
ALTER COLUMN request_time SET NOT NULL;

ALTER TABLE database_connection_info
DROP COLUMN dbName;

ALTER TABLE database_connection_info
ADD COLUMN dbname TEXT NOT NULL;

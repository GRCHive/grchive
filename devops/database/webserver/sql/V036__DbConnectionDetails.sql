ALTER TABLE database_connection_info
DROP COLUMN connection_string;

ALTER TABLE database_connection_info
ADD COLUMN host TEXT NOT NULL,
ADD COLUMN port INTEGER NOT NULL,
ADD COLUMN dbName TEXT NOT NULL,
ADD COLUMN parameters JSONB;

DELETE FROM database_connection_info;

DELETE FROM database_refresh;
DROP TABLE database_columns;

ALTER TABLE database_tables
ADD COLUMN columns JSONB;

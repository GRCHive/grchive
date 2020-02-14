ALTER TABLE database_sql_metadata
DROP COLUMN name;

ALTER TABLE database_sql_metadata
ADD COLUMN name VARCHAR(255) NOT NULL;

CREATE UNIQUE INDEX ON database_sql_metadata(db_id, name);

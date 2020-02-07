ALTER TABLE supported_databases
ADD COLUMN has_sql_support BOOLEAN NOT NULL DEFAULT true;

UPDATE supported_databases
SET has_sql_support = false
WHERE name = 'Other';

ALTER TABLE supported_databases
ALTER COLUMN has_sql_support
DROP DEFAULT;

CREATE INDEX ON database_resources(type_id);

UPDATE supported_databases
SET has_sql_support = true
WHERE name != 'Other';

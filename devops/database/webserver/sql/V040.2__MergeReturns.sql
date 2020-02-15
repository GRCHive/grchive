DROP TABLE database_function_returns;

ALTER TABLE database_functions
ADD COLUMN ret_type TEXT;

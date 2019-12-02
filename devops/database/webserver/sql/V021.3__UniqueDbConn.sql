ALTER TABLE database_connection_info
ADD CONSTRAINT unique_db_conn_info UNIQUE(org_id, db_id);

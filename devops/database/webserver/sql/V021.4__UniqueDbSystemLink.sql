ALTER TABLE database_system_link
ADD CONSTRAINT unique_org_db_system UNIQUE(org_id, db_id, system_id);

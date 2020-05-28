ALTER TABLE shell_scripts
ALTER COLUMN bucket_id DROP NOT NULL;

ALTER TABLE shell_scripts
ALTER COLUMN storage_id DROP NOT NULL;

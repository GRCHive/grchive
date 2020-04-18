ALTER TABLE script_runs
ADD COLUMN requires_build BOOLEAN DEFAULT true;

UPDATE script_runs
SET requires_build = false
WHERE build_start_time IS NULL;

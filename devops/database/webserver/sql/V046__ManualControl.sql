ALTER TABLE process_flow_controls
ADD COLUMN is_manual BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE process_flow_controls
ALTER COLUMN is_manual DROP DEFAULT;

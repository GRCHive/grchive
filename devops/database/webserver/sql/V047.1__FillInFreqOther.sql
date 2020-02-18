UPDATE process_flow_controls
SET freq_other = '';

ALTER TABLE process_flow_controls
ALTER COLUMN freq_other SET NOT NULL;

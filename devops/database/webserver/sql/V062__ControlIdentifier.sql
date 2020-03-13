ALTER TABLE process_flow_controls
ADD COLUMN identifier VARCHAR(256);

UPDATE process_flow_controls AS c
SET identifier = c.name;

ALTER TABLE process_flow_controls
ALTER COLUMN identifier SET NOT NULL;

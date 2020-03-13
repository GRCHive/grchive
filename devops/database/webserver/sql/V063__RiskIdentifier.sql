ALTER TABLE process_flow_risks
ADD COLUMN identifier VARCHAR(256);

UPDATE process_flow_risks AS c
SET identifier = c.name;

ALTER TABLE process_flow_risks
ALTER COLUMN identifier SET NOT NULL;

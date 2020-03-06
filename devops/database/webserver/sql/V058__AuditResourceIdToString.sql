ALTER TABLE global_audit_event_history
ADD COLUMN resource_id_str VARCHAR(128);

INSERT INTO global_audit_event_history (resource_id_str)
SELECT CAST(resource_id AS VARCHAR(128))
FROM global_audit_event_history;

ALTER TABLE global_audit_event_history
ALTER COLUMN resource_id_str SET NOT NULL;

ALTER TABLE global_audit_event_history
DROP COLUMN resource_id;

ALTER TABLE global_audit_event_history
RENAME COLUMN resource_id_str TO resource_id;

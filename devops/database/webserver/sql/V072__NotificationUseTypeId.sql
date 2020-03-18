ALTER TABLE notifications
RENAME COLUMN subject TO subject_type;

ALTER TABLE notifications
RENAME COLUMN object TO object_type;

ALTER TABLE notifications
RENAME COLUMN indirect_object TO indirect_object_type;

ALTER TABLE notifications
DROP COLUMN subject_data;

ALTER TABLE notifications
DROP COLUMN object_data;

ALTER TABLE notifications
DROP COLUMN indirect_object_data;

ALTER TABLE notifications
ADD COLUMN subject_id BIGINT NOT NULL;

ALTER TABLE notifications
ADD COLUMN object_id BIGINT NOT NULL;

ALTER TABLE notifications
ADD COLUMN indirect_object_id BIGINT NOT NULL;

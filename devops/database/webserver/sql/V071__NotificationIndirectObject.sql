ALTER TABLE notifications
ADD COLUMN indirect_object TEXT NOT NULL;

ALTER TABLE notifications
ADD COLUMN indirect_object_data JSONB NOT NULL;

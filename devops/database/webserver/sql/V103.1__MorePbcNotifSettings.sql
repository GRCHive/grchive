CREATE UNIQUE INDEX ON org_pbc_notification_cadence_settings(org_id, days_before_due);

ALTER TABLE org_pbc_notification_cadence_settings 
ADD COLUMN send_to_assignee BOOLEAN DEFAULT true;

ALTER TABLE org_pbc_notification_cadence_settings 
ADD COLUMN send_to_requester BOOLEAN DEFAULT false;

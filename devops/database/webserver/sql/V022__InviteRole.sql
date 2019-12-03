ALTER TABLE invitation_codes
ADD COLUMN role_id BIGINT NOT NULL;

ALTER TABLE invitation_codes
ADD CONSTRAINT role_org_foreign_key FOREIGN KEY (role_id, from_org_id) REFERENCES organization_available_roles (id, org_id);

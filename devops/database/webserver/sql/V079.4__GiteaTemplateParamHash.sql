ALTER TABLE organization_gitea_repository_template
ADD COLUMN template256sum VARCHAR(64) NOT NULL DEFAULT '';

ALTER TABLE organization_gitea_repository_template
ALTER COLUMN template256sum DROP DEFAULT;

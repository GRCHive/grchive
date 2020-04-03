-- Table to track which version of the template the organization's repository is currently at.
CREATE TABLE organization_gitea_repository_template (
    org_id INTEGER NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    sha256sum VARCHAR(64) NOT NULL
);

CREATE INDEX ON organization_gitea_repository_template (sha256sum);

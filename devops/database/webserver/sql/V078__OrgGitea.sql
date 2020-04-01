CREATE TABLE org_gitea_link (
    org_id INTEGER NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    gitea_organization TEXT NOT NULL,
    gitea_repository TEXT NOT NULL,
    gitea_user TEXT NOT NULL,
    gitea_access_token_vault_secret TEXT NOT NULL
);

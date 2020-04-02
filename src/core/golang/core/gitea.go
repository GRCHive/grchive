package core

type LinkedGiteaRepository struct {
	OrgId                     int32  `db:"org_id"`
	GiteaOrg                  string `db:"gitea_organization"`
	GiteaRepo                 string `db:"gitea_repository"`
	GiteaUser                 string `db:"gitea_user"`
	GiteaAccessTokenVaultPath string `db:"gitea_access_token_vault_secret"`
}

package database

import (
	"gitlab.com/grchive/grchive/core"
)

func LinkOrganizationToGitea(orgId int32, giteaOrg string, giteaRepository string, giteaUsername string, giteaAccessVaultSecret string) error {
	tx := CreateTx()
	_, err := tx.Exec(`
		INSERT INTO org_gitea_link (org_id, gitea_organization, gitea_repository, gitea_user, gitea_access_token_vault_secret)
		VALUES ($1, $2, $3, $4, $5)
	`, orgId, giteaOrg, giteaRepository, giteaUsername, giteaAccessVaultSecret)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetLinkedGiteaRepository(orgId int32) (*core.LinkedGiteaRepository, error) {
	repo := core.LinkedGiteaRepository{}
	err := dbConn.Get(&repo, `
		SELECT *
		FROM org_gitea_link
		WHERE org_id = $1
	`, orgId)
	return &repo, err
}

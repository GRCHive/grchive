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

func GetGiteaTemplateHashForOrg(orgId int32) (string, string, error) {
	rows, err := dbConn.Queryx(`
		SELECT sha256sum, template256sum
		FROM organization_gitea_repository_template
		WHERE org_id = $1
	`, orgId)

	if err != nil {
		return "", "", err
	}

	defer rows.Close()

	if !rows.Next() {
		return "", "", nil
	}

	hash := ""
	templateHash := ""
	err = rows.Scan(&hash, &templateHash)
	if err != nil {
		return "", "", err
	}

	return hash, templateHash, nil
}

func StoreGiteaTemplateHashForOrg(orgId int32, hash string, templateHash string) error {
	tx := CreateTx()
	_, err := tx.Exec(`
		INSERT INTO organization_gitea_repository_template (org_id, sha256sum, template256sum)
		VALUES ($1, $2, $3)
		ON CONFLICT (org_id) DO UPDATE SET
			sha256sum = EXCLUDED.sha256sum,
			template256sum = EXCLUDED.template256sum
	`, orgId, hash, templateHash)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

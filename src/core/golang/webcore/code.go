package webcore

import (
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
)

func GetManagedCodeFromGitea(codeId int64, orgId int32, role *core.Role) (string, error) {
	grcRepo, err := database.GetLinkedGiteaRepository(orgId)
	if err != nil {
		return "", err
	}

	code, err := database.GetCode(codeId, orgId, role)
	if err != nil {
		return "", err
	}

	content, _, err := gitea.GlobalGiteaApi.RepositoryGetFile(gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}, code.GitPath, code.GitHash)
	return content, err
}

func StoreManagedCodeToGitea(code *core.ManagedCode, script string, role *core.Role) error {
	grcRepo, err := database.GetLinkedGiteaRepository(code.OrgId)
	if err != nil {
		return err
	}

	err = UpdateGiteaRepositoryTemplate(code.OrgId)
	if err != nil {
		return err
	}

	giteaRepo := gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}

	// We need to check if the file exists or not in Gitea irrespective of what our database thinks
	// since it's possible for our code to get into a state where the file exists in Gitea but we do not
	// have it tracked in the database.
	_, sha, err := gitea.GlobalGiteaApi.RepositoryGetFile(giteaRepo, code.GitPath, "master")

	// Assume a non-nil error means that the file doesn't exist. If it turns out we're wrong,
	// we should get another error when we try to create a new file.
	if err != nil {
		code.GitHash, code.GiteaFileSha, err = gitea.GlobalGiteaApi.RepositoryCreateFile(
			giteaRepo,
			code.GitPath,
			gitea.GiteaCreateFileOptions{
				Content: script,
			},
		)
	} else {
		code.GitHash, code.GiteaFileSha, err = gitea.GlobalGiteaApi.RepositoryUpdateFile(
			giteaRepo,
			code.GitPath,
			gitea.GiteaCreateFileOptions{
				Content: script,
			},
			sha,
		)
	}

	if err != nil {
		return err
	}

	// Store information in the database after since it'll be easier to handle cases where
	// files are in Gitea but not in the database than if the file were in the database but not in Gitea.
	return database.InsertManagedCode(code, role)
}

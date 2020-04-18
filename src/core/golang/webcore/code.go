package webcore

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"gopkg.in/yaml.v2"
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

func DeleteManagedCodeFromGitea(code *core.ManagedCode) error {
	grcRepo, err := database.GetLinkedGiteaRepository(code.OrgId)
	if err != nil {
		return err
	}

	giteaRepo := gitea.GiteaRepository{
		Owner: grcRepo.GiteaOrg,
		Name:  grcRepo.GiteaRepo,
	}

	err = gitea.GlobalGiteaApi.RepositoryDeleteFile(giteaRepo, code.GitPath, gitea.GiteaDeleteFileOptions{
		Message: fmt.Sprintf("Delete file - %s", code.GitPath),
		Sha:     code.GiteaFileSha,
	})

	return err
}

func StoreManagedCodeToGitea(code *core.ManagedCode, script string, role *core.Role, msg string) error {
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
				Message: msg,
			},
		)
	} else {
		code.GitHash, code.GiteaFileSha, err = gitea.GlobalGiteaApi.RepositoryUpdateFile(
			giteaRepo,
			code.GitPath,
			gitea.GiteaCreateFileOptions{
				Content: script,
				Message: msg,
			},
			sha,
		)
	}

	if err != nil {
		return err
	}

	// Store information in the database after since it'll be easier to handle cases where
	// files are in Gitea but not in the database than if the file were in the database but not in Gitea.
	if role != nil {
		code.UserId = role.UserId
		return database.InsertManagedCode(code, role)
	}
	return nil
}

func GenerateScriptMetadataYaml(params []*core.CodeParameter, clientDataId []int64) (string, error) {
	type ScriptMetadataData struct {
		Params       []*core.CodeParameter `yaml:"params"`
		ClientDataId []int64               `yaml:"clientDataId"`
	}

	metadata := ScriptMetadataData{
		Params:       params,
		ClientDataId: clientDataId,
	}

	data, err := yaml.Marshal(&metadata)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

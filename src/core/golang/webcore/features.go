package webcore

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	"gitlab.com/grchive/grchive/gitea_api"
	"gitlab.com/grchive/grchive/vault_api"
)

func EnableAutomationFeature(orgId int32) error {
	grchiveOrg, err := database.FindOrganizationFromId(orgId)
	if err != nil {
		return err
	}

	// This functioln needs to setup organization specific things
	// in Gitea and Drone CI.
	// 	1) Create an organization specific user for us to assume the identity of.
	// 	2) Create an organization and repository for holding all the
	// 	   org's Kotlin code in Gitea. Note that due the limitiations of the Gitea
	// 	   API I *think* we have to create the repository as the user and then transfer
	// 	   the repository to the organization. Doesn't really matter but might as well
	// 	   organize.
	//  3) Put in default template code for the project.
	//  4) Enable the repository in Drone CI.
	pw, err := core.RandomHexString(32)
	if err != nil {
		return err
	}

	// Need to meet Gitea password complexity requirements.
	pw = pw + "A!"

	user := gitea.GiteaUser{
		Username: "grchive-" + grchiveOrg.OktaGroupName,
		Password: pw,
		Email:    fmt.Sprintf("gitea+%s@grchive.com", grchiveOrg.OktaGroupName),
		FullName: grchiveOrg.Name,
	}

	err = gitea.GlobalGiteaApi.AdminCreateUser(user)
	if err != nil {
		return err
	}

	token, err := gitea.GlobalGiteaApi.UserCreateAccessToken(user, "grchive-webserver-access-token")
	if err != nil {
		return err
	}

	giteaOrg := gitea.GiteaOrganization{
		Username: grchiveOrg.OktaGroupName,
		FullName: grchiveOrg.Name,
	}

	err = gitea.GlobalGiteaApi.AdminCreateOrganization(user, giteaOrg)
	if err != nil {
		return err
	}

	repository := gitea.GiteaRepository{
		Name: fmt.Sprintf("grchive-automation-%s", grchiveOrg.OktaGroupName),
	}

	err = gitea.GlobalGiteaApi.RepositoryCreate(*token, repository)
	if err != nil {
		return err
	}

	err = gitea.GlobalGiteaApi.RepositoryTransfer(user, giteaOrg, repository)
	if err != nil {
		return err
	}

	userTokenVaultPath := fmt.Sprintf("secret/webserver/gitea/tokens/%s", user.Username)
	err = vault.StoreSecret(userTokenVaultPath, map[string]string{
		"name":  token.Name,
		"token": token.Token,
	}, -1)
	if err != nil {
		return err
	}

	err = database.LinkOrganizationToGitea(orgId, giteaOrg.Username, repository.Name, user.Username, userTokenVaultPath)
	if err != nil {
		return err
	}

	return nil
}

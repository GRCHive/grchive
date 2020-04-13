package webcore

import (
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"gitlab.com/grchive/grchive/database"
	drone "gitlab.com/grchive/grchive/drone_api"
	gitea "gitlab.com/grchive/grchive/gitea_api"
	vault "gitlab.com/grchive/grchive/vault_api"
)

func enableAutomationGitea(grchiveOrg *core.Organization) error {
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

	core.Debug(" -- Gitea Admin Create User")
	err = gitea.GlobalGiteaApi.AdminCreateUser(user)
	if err != nil {
		return err
	}

	core.Debug(" -- Gitea Create Access Token")
	token, err := gitea.GlobalGiteaApi.UserCreateAccessToken(user, "grchive-webserver-access-token")
	if err != nil {
		return err
	}

	repository := gitea.GiteaRepository{
		Name:  fmt.Sprintf("grchive-automation-%s", grchiveOrg.OktaGroupName),
		Owner: user.Username,
	}

	core.Debug(" -- Gitea Repository Create")
	err = gitea.GlobalGiteaApi.RepositoryCreate(*token, repository)
	if err != nil {
		return err
	}

	core.Debug(" -- Gitea Repository Transfer")
	err = gitea.GlobalGiteaApi.RepositoryTransfer(user, gitea.GiteaOrganization{
		Username: core.EnvConfig.Gitea.GlobalOrg,
	}, &repository)
	if err != nil {
		return err
	}

	// Add the admin user as a collaborator on the repository so that we can access it
	// from Drone CI.
	// TODO: Move admin username to config.
	core.Debug(" -- Gitea Repository Add Collaborator")
	err = gitea.GlobalGiteaApi.RepositoryAddCollaborator(repository, gitea.GiteaUser{
		Username: "grchive-gitea-admin",
	})

	if err != nil {
		return err
	}

	core.Debug(" -- Gitea Vault Store Secret")
	userTokenVaultPath := fmt.Sprintf("secret/webserver/gitea/tokens/%s", user.Username)
	err = vault.StoreSecret(userTokenVaultPath, map[string]string{
		"name":  token.Name,
		"token": token.Token,
	}, -1)
	if err != nil {
		return err
	}

	core.Debug(" -- Gitea Link Organization to Gitea")
	err = database.LinkOrganizationToGitea(grchiveOrg.Id, repository.Owner, repository.Name, user.Username, userTokenVaultPath)
	if err != nil {
		return err
	}

	return nil
}

func enableAutomationDrone(org *core.Organization) error {
	repo, err := database.GetLinkedGiteaRepository(org.Id)
	if err != nil {
		return err
	}

	// Need to force a sync here or else the repository won't
	// show up in Drone and thus enabling it will do nothing.
	err = drone.GlobalDroneApi.RepoSync()
	if err != nil {
		return err
	}

	err = drone.GlobalDroneApi.RepoEnable(repo.GiteaOrg, repo.GiteaRepo)
	if err != nil {
		return err
	}
	return nil
}

func EnableAutomationFeature(orgId int32) error {
	grchiveOrg, err := database.FindOrganizationFromId(orgId)
	if err != nil {
		return err
	}

	// This function needs to setup organization specific things
	// in Gitea and Drone CI.
	// 	1) Create an organization specific user for us to assume the identity of.
	// 	2) Create an organization and repository for holding all the
	// 	   org's Kotlin code in Gitea. Note that due the limitiations of the Gitea
	// 	   API I *think* we have to create the repository as the user and then transfer
	// 	   the repository to the organization. Doesn't really matter but might as well
	// 	   organize.
	//  3) Put in default template code for the project.
	//  4) Enable the repository in Drone CI.
	err = enableAutomationGitea(grchiveOrg)
	if err != nil {
		return err
	}

	err = enableAutomationDrone(grchiveOrg)
	if err != nil {
		return err
	}

	err = UpdateGiteaRepositoryTemplate(grchiveOrg.Id)
	if err != nil {
		return err
	}

	return nil
}
package gitea

import (
	"fmt"
)

const OrgCreateRepoEndpoint = "/orgs/%s/repos"

func (r *RealGiteaApi) CreateRepositoryForOrganization(repo GiteaRepository, org GiteaOrganization) error {
	_, _, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(OrgCreateRepoEndpoint, org.Username),
		r.cfg.Token,
		map[string]interface{}{
			"name":    repo.Name,
			"private": true,
		},
	)
	return err
}

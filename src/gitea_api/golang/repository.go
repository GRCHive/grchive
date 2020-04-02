package gitea

import (
	"fmt"
)

const CreateRepoEndpoint = "/user/repos"
const TransferRepoEndpoint = "/repos/%s/%s/transfer"
const AddCollabEndpoint = "/repos/%s/%s/collaborators/%s"

func (r *RealGiteaApi) RepositoryCreate(token GiteaToken, repo GiteaRepository) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		CreateRepoEndpoint,
		token.Token,
		map[string]interface{}{
			"name": repo.Name,
		},
	)
	return err
}

func (r *RealGiteaApi) RepositoryTransfer(from GiteaUserlike, to GiteaUserlike, repo GiteaRepository) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(TransferRepoEndpoint, from.GetUsername(), repo.Name),
		r.cfg.Token,
		map[string]interface{}{
			"new_owner": to.GetUsername(),
		},
	)
	return err
}

func (r *RealGiteaApi) RepositoryAddCollaborator(repo GiteaRepository, owner GiteaUserlike, collab GiteaUserlike) error {
	_, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(AddCollabEndpoint,
			owner.GetUsername(),
			repo.Name,
			collab.GetUsername(),
		),
		r.cfg.Token,
		map[string]interface{}{},
	)
	return err
}

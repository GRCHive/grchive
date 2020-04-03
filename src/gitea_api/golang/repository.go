package gitea

import (
	"fmt"
)

const CreateRepoEndpoint = "/user/repos"
const TransferRepoEndpoint = "/repos/%s/%s/transfer"
const AddCollabEndpoint = "/repos/%s/%s/collaborators/%s"
const FileContentEndpoint = "/repos/%s/%s/contents/%s"

func (r *RealGiteaApi) RepositoryCreate(token GiteaToken, repo GiteaRepository) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		CreateRepoEndpoint,
		token.Token,
		map[string]interface{}{
			"auto_init":      true,
			"default_branch": "master",
			"name":           repo.Name,
			"private":        true,
		},
	)
	return err
}

func (r *RealGiteaApi) RepositoryTransfer(from GiteaUserlike, to GiteaUserlike, repo *GiteaRepository) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(TransferRepoEndpoint, from.GetUsername(), repo.Name),
		r.cfg.Token,
		map[string]interface{}{
			"new_owner": to.GetUsername(),
		},
	)

	if err != nil {
		repo.Owner = to.GetUsername()
	}

	return err
}

func (r *RealGiteaApi) RepositoryAddCollaborator(repo GiteaRepository, collab GiteaUserlike) error {
	_, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(AddCollabEndpoint,
			repo.Owner,
			repo.Name,
			collab.GetUsername(),
		),
		r.cfg.Token,
		map[string]interface{}{},
	)
	return err
}

func (r *RealGiteaApi) RepositoryCreateFile(repo GiteaRepository, path string, content string) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		nil,
	)
	return err
}

func (r *RealGiteaApi) RepositoryUpdateFile(repo GiteaRepository, path string, content string) error {
	_, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		nil,
	)
	return err
}

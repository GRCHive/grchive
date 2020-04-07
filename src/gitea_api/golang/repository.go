package gitea

import (
	"encoding/base64"
	"encoding/json"
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

func getCommitFileShaFromResponse(resp map[string]*json.RawMessage) (string, string, error) {
	commitData := map[string]interface{}{}
	err := json.Unmarshal(*resp["commit"], &commitData)
	if err != nil {
		return "", "", err
	}

	contentData := map[string]interface{}{}
	err = json.Unmarshal(*resp["content"], &contentData)
	if err != nil {
		return "", "", err
	}

	return commitData["sha"].(string), contentData["sha"].(string), nil
}

func (r *RealGiteaApi) RepositoryCreateFile(repo GiteaRepository, path string, content string) (string, string, error) {
	base64Data := base64.StdEncoding.EncodeToString([]byte(content))
	resp, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		map[string]interface{}{
			"content": base64Data,
		},
	)

	if err != nil {
		return "", "", err
	}

	return getCommitFileShaFromResponse(resp)
}

func (r *RealGiteaApi) RepositoryUpdateFile(repo GiteaRepository, path string, content string, sha string) (string, string, error) {
	base64Data := base64.StdEncoding.EncodeToString([]byte(content))
	resp, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		map[string]interface{}{
			"content": base64Data,
			"sha":     sha,
		},
	)

	if err != nil {
		return "", "", err
	}

	return getCommitFileShaFromResponse(resp)
}

func (r *RealGiteaApi) RepositoryGetFile(repo GiteaRepository, path string, ref string) (string, string, error) {
	resp, err := r.sendGiteaRequestWithToken(
		"GET",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		)+"?ref="+ref,
		r.cfg.Token,
		nil,
	)

	if err != nil {
		return "", "", err
	}

	content := ""
	sha := ""

	err = json.Unmarshal(*resp["content"], &content)
	if err != nil {
		return "", "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(*resp["sha"], &sha)
	if err != nil {
		return "", "", err
	}

	return string(decoded), sha, err
}

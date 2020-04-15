package gitea

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const CreateRepoEndpoint = "/user/repos"
const TransferRepoEndpoint = "/repos/%s/%s/transfer"
const AddCollabEndpoint = "/repos/%s/%s/collaborators/%s"
const FileContentEndpoint = "/repos/%s/%s/contents/%s"
const GitSingleRefEndpoint = "/repos/%s/%s/git/refs/%s"

func (r *RealGiteaApi) RepositoryCreate(token GiteaToken, repo GiteaRepository) error {
	_, _, err := r.sendGiteaRequestWithToken(
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
	_, _, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(TransferRepoEndpoint, from.GetUsername(), repo.Name),
		r.cfg.Token,
		map[string]interface{}{
			"new_owner": to.GetUsername(),
		},
	)

	if err == nil {
		repo.Owner = to.GetUsername()
	}

	return err
}

func (r *RealGiteaApi) RepositoryAddCollaborator(repo GiteaRepository, collab GiteaUserlike) error {
	sc, _, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(AddCollabEndpoint,
			repo.Owner,
			repo.Name,
			collab.GetUsername(),
		),
		r.cfg.Token,
		map[string]interface{}{},
	)

	if sc == http.StatusNoContent {
		return nil
	}

	return err
}

func getCommitFileShaFromResponse(resp json.RawMessage) (string, string, error) {
	dictResponse := map[string]*json.RawMessage{}
	err := json.Unmarshal(resp, &dictResponse)
	if err != nil {
		return "", "", err
	}

	commitData := map[string]interface{}{}
	err = json.Unmarshal(*dictResponse["commit"], &commitData)
	if err != nil {
		return "", "", err
	}

	contentData := map[string]interface{}{}
	err = json.Unmarshal(*dictResponse["content"], &contentData)
	if err != nil {
		return "", "", err
	}

	return commitData["sha"].(string), contentData["sha"].(string), nil
}

func (r *RealGiteaApi) RepositoryCreateFile(repo GiteaRepository, path string, opts GiteaCreateFileOptions) (string, string, error) {
	_, resp, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		opts.PrepareApiBody(),
	)

	if err != nil {
		return "", "", err
	}

	return getCommitFileShaFromResponse(resp)
}

func (r *RealGiteaApi) RepositoryUpdateFile(repo GiteaRepository, path string, opts GiteaCreateFileOptions, sha string) (string, string, error) {
	body := opts.PrepareApiBody()
	body["sha"] = sha

	_, resp, err := r.sendGiteaRequestWithToken(
		"PUT",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		body,
	)

	if err != nil {
		return "", "", err
	}

	return getCommitFileShaFromResponse(resp)
}

func (r *RealGiteaApi) RepositoryGetFile(repo GiteaRepository, path string, ref string) (string, string, error) {
	_, resp, err := r.sendGiteaRequestWithToken(
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

	dictResponse := map[string]*json.RawMessage{}
	err = json.Unmarshal(resp, &dictResponse)
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(*dictResponse["content"], &content)
	if err != nil {
		return "", "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", "", err
	}

	err = json.Unmarshal(*dictResponse["sha"], &sha)
	if err != nil {
		return "", "", err
	}

	return string(decoded), sha, err
}

func (r *RealGiteaApi) RepositoryDeleteFile(repo GiteaRepository, path string, opts GiteaDeleteFileOptions) error {
	_, _, err := r.sendGiteaRequestWithToken(
		"DELETE",
		fmt.Sprintf(FileContentEndpoint,
			repo.Owner,
			repo.Name,
			path,
		),
		r.cfg.Token,
		opts,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RealGiteaApi) RepositoryGitGetRefSha(repo GiteaRepository, ref string) (string, error) {
	_, resp, err := r.sendGiteaRequestWithToken(
		"GET",
		fmt.Sprintf(GitSingleRefEndpoint,
			repo.Owner,
			repo.Name,
			strings.TrimPrefix(ref, "refs/"),
		),
		r.cfg.Token,
		nil,
	)

	if err != nil {
		return "", err
	}

	arrResponse := []json.RawMessage{}
	err = json.Unmarshal(resp, &arrResponse)
	if err != nil {
		return "", err
	}

	if len(arrResponse) == 0 {
		return "", errors.New("Failed to find ref.")
	}

	refData := map[string]json.RawMessage{}
	err = json.Unmarshal(arrResponse[0], &refData)
	if err != nil {
		return "", err
	}

	objData := map[string]string{}
	err = json.Unmarshal(refData["object"], &objData)
	if err != nil {
		return "", err
	}

	return objData["sha"], nil
}

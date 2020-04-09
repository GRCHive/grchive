package gitea

import (
	"encoding/json"
	"fmt"
)

const UserTokenEndpoint = "/users/%s/tokens"

func (r *RealGiteaApi) UserCreateAccessToken(user GiteaUser, tokenName string) (*GiteaToken, error) {
	_, data, err := r.sendGiteaRequestWithUserAuth(
		"POST",
		fmt.Sprintf(UserTokenEndpoint, user.Username),
		user,
		map[string]interface{}{
			"name": tokenName,
		},
	)

	if err != nil {
		return nil, err
	}

	// The docs (https://try.gitea.io/api/swagger#/user/userCreateToken) say that
	// the results should come back in the header but getting the data from the body
	// works fine so...
	token := GiteaToken{}

	err = json.Unmarshal(*data["name"], &token.Name)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*data["sha1"], &token.Token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

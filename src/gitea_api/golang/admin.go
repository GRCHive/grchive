package gitea

import (
	"fmt"
)

const AdminCreateUserEndpoint string = "/admin/users"
const AdminCreateOrgEndpoint string = "/admin/users/%s/orgs"

func (r *RealGiteaApi) AdminCreateUser(user GiteaUser) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		AdminCreateUserEndpoint,
		r.cfg.Token,
		map[string]interface{}{
			"email":                user.Email,
			"full_name":            user.FullName,
			"login_name":           user.Username,
			"must_change_password": false,
			"password":             user.Password,
			"send_notify":          false,
			"username":             user.Username,
		},
	)
	return err
}

func (r *RealGiteaApi) AdminCreateOrganization(user GiteaUser, org GiteaOrganization) error {
	_, err := r.sendGiteaRequestWithToken(
		"POST",
		fmt.Sprintf(AdminCreateOrgEndpoint, user.Username),
		r.cfg.Token,
		map[string]interface{}{
			"full_name": org.FullName,
			"username":  org.Username,
		},
	)
	return err
}

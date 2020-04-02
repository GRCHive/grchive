package drone

import (
	"fmt"
)

const RepoEnableEndpoint = "/api/repos/%s/%s"
const UserRepoEndpoint = "/api/user/repos"

func (d *RealDroneApi) RepoEnable(owner string, repo string) error {
	_, err := d.sendDroneRequest(
		"POST",
		fmt.Sprintf(RepoEnableEndpoint, owner, repo),
		d.cfg.Token,
		nil,
	)
	return err
}

func (d *RealDroneApi) RepoSync() error {
	_, err := d.sendDroneRequest(
		"POST",
		UserRepoEndpoint,
		d.cfg.Token,
		nil,
	)
	return err
}

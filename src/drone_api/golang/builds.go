package drone

import (
	"fmt"
	"net/url"
	"strings"
)

const BuildCreateEndpoint = "/api/repos/%s/%s/builds"

func (d *RealDroneApi) BuildCreate(owner string, repo string, params map[string]string) error {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(BuildCreateEndpoint, owner, repo))
	builder.WriteString("?branch=master")
	for k, v := range params {
		builder.WriteString(fmt.Sprintf("&%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}

	_, err := d.sendDroneRequest(
		"POST",
		builder.String(),
		d.cfg.Token,
		nil,
	)
	return err
}

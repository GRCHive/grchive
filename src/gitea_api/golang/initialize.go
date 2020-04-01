package gitea

import (
	"net/http"
)

func (r *RealGiteaApi) MustInitialize(cfg GiteaConfig) {
	r.cfg = cfg
	r.client = &http.Client{}
}

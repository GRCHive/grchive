package webcore

import (
	"gitlab.com/grchive/grchive/core"
	"net/http"
)

type HttpHandler = func(http.ResponseWriter, *http.Request)

func CreateACLCheckPermissionHandler(f HttpHandler, access ...core.ResourceAccessBundle) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		role, err := FindRoleInContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		for _, a := range access {
			if !role.Permissions.HasAccess(a.Resource, a.Access) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		f(w, r)
	}
}

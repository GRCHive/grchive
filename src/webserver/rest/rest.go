package rest

import (
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
)

func RegisterPaths(r *mux.Router) {
	r.HandleFunc(core.CreateGetStartedUrl(), postGettingStartedInterest).Methods("POST")
}

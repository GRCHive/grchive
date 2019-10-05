package webcore

import (
	"errors"
	"github.com/gorilla/mux"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
)

func GetOrganizationFromRequestUrl(r *http.Request) (*core.Organization, error) {
	urlRouteVars := mux.Vars(r)
	orgGroupName, ok := urlRouteVars[core.DashboardOrgOrgQueryId]
	if !ok {
		return nil, errors.New("No valid organization in request URL")
	}

	org, err := database.FindOrganizationFromGroupName(orgGroupName)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func GetUserEmailFromRequestUrl(r *http.Request) (string, error) {
	urlRouteVars := mux.Vars(r)
	email, ok := urlRouteVars[core.DashboardUserQueryId]
	if !ok {
		return "", errors.New("No email in request URL")
	}

	return email, nil
}
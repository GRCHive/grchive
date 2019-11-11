package webcore

import (
	"errors"
	"gitlab.com/b3h47pte/audit-stuff/core"
	"gitlab.com/b3h47pte/audit-stuff/database"
	"net/http"
)

func ObtainOrganizationDefaultRole(orgId int32) (*core.Role, error) {
	role, err := database.FindDefaultRoleForOrg(orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	if role != nil {
		return role, nil
	}

	// No default permissions was found which probably means this org was just created
	// so create a default default permissions which is full admin access to everybody.
	defaultRole := core.Role{
		Permissions: core.CreateAllAccessPermission(),
	}
	defaultMetadata := core.CreateDefaultRoleMetadata(orgId)

	err = database.InsertOrgRole(&defaultMetadata, &defaultRole, core.ServerRole)
	if err != nil {
		return nil, err
	}

	return &defaultRole, nil
}

func ObtainAPIKeyRole(key *core.ApiKey, orgId int32) (*core.Role, error) {
	role, err := database.FindUserRoleForOrg(key.UserId, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	if role != nil {
		return role, nil
	}

	return GrantAPIKeyDefaultRole(key, orgId)
}

func GrantAPIKeyDefaultRole(key *core.ApiKey, orgId int32) (*core.Role, error) {
	// At this point we know that the user doesn't have a set permissions yet so we need
	// to give the user default access controls (as specified by the org).
	// We only give these default access controls to users belonging to the organization.
	user, err := database.FindUserFromId(key.UserId)
	if err != nil {
		return nil, err
	}

	accessibleOrgIds, err := database.FindAccessibleOrganizationsForUser(user)
	if err != nil {
		return nil, err
	}

	if core.LinearSearchInt32Slice(accessibleOrgIds, orgId) == core.SearchNotFound {
		return nil, errors.New("User does not have access.")
	}

	defaultRole, err := ObtainOrganizationDefaultRole(orgId)
	if err != nil {
		return nil, err
	}

	err = database.InsertUserRoleForOrg(user.Id, orgId, defaultRole, core.ServerRole)
	if err != nil {
		return nil, err
	}

	return defaultRole, nil
}

func GetCurrentRequestRole(r *http.Request, orgId int32) (*core.Role, error) {
	key, err := FindApiKeyInContext(r.Context())
	if err != nil {
		return nil, err
	}
	return ObtainAPIKeyRole(key, orgId)
}

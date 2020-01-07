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
		Permissions: core.CreateViewOnlyAccessPermission(),
	}
	defaultMetadata := core.CreateDefaultRoleMetadata(orgId)

	err = database.InsertOrgRole(&defaultMetadata, &defaultRole, core.ServerRole)
	if database.IsDuplicateDBEntry(err) {
		// Something happened we should actually be able to find this thing...try again!
		return ObtainOrganizationDefaultRole(orgId)
	} else if err != nil {
		return nil, err
	}

	return &defaultRole, nil
}

func ObtainOrganizationAdminRole(orgId int32) (*core.Role, error) {
	role, err := database.FindAdminRoleForOrg(orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	if role != nil {
		return role, nil
	}

	// No default permissions was found which probably means this org was just created
	// so create a default default permissions which is full admin access to everybody.
	adminRole := core.Role{
		Permissions: core.CreateAllAccessPermission(),
	}
	adminMetadata := core.CreateAdminRoleMetadata(orgId)

	err = database.InsertOrgRole(&adminMetadata, &adminRole, core.ServerRole)
	if database.IsDuplicateDBEntry(err) {
		// Something happened we should actually be able to find this thing...try again!
		return ObtainOrganizationAdminRole(orgId)
	} else if err != nil {
		return nil, err
	}

	return &adminRole, nil
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

	accessibleOrgIds, err := database.FindAccessibleOrganizationIdsForUser(user)
	if err != nil {
		return nil, err
	}

	if core.LinearSearchInt32Slice(accessibleOrgIds, orgId) == core.SearchNotFound {
		return nil, errors.New("User does not have access.")
	}

	// Obtain both the default role and admin role to make sure the
	// organization's default roles are initialized.
	defaultRole, err := ObtainOrganizationDefaultRole(orgId)
	if err != nil {
		return nil, err
	}

	adminRole, err := ObtainOrganizationAdminRole(orgId)
	if err != nil {
		return nil, err
	}

	adminUserIds, err := database.FindUserIdsWithRole(adminRole.Id, orgId, core.ServerRole)
	if err != nil {
		return nil, err
	}

	// If this user is the first in their organization to login, they are the de-facto admin.
	if len(adminUserIds) == 0 {
		err = database.InsertUserRoleForOrg(user.Id, orgId, adminRole.Id, core.ServerRole)
	} else {
		err = database.InsertUserRoleForOrg(user.Id, orgId, defaultRole.Id, core.ServerRole)
	}

	// It's ok if there's a duplicate because it means we've added the role already. OK!
	if err != nil && !database.IsDuplicateDBEntry(err) {
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

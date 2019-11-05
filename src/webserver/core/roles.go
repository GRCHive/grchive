package core

type AccessType struct {
	// View: Can see the thing being granted access to.
	CanView bool
	// Edit: Can change the thing being granted access to.
	CanEdit bool
	// Manage: Can add/delete the thing being granted access to.
	CanManage bool
}

type ResourceType int

const (
	// Managing is merely for deleting/creating new process flows.
	OrgRoles ResourceType = iota
	// In the case of manage, allows giving/revoking a role from a user.
	// In this case, the editing the process flow includes adding/creating nodes.
	ProcessFlows
	Controls
	ControlDocumentation
	Risks
)

type PermissionsMap map[ResourceType]AccessType

type RoleMetadata struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IsDefault   bool   `db:"is_default_role"`
	OrgId       int32  `db:"org_id"`
}

type Role struct {
	Id          int64          `db:"id"`
	Permissions PermissionsMap `db:"permissions"`
}

func CreateOwnerAccessType() AccessType {
	return AccessType{
		CanView:   true,
		CanEdit:   true,
		CanManage: true,
	}
}

func CreateAllAccessPermission() PermissionsMap {
	return PermissionsMap{
		OrgRoles:             CreateOwnerAccessType(),
		ProcessFlows:         CreateOwnerAccessType(),
		Controls:             CreateOwnerAccessType(),
		ControlDocumentation: CreateOwnerAccessType(),
		Risks:                CreateOwnerAccessType(),
	}
}

func CreateDefaultRoleMetadata(orgId int32) RoleMetadata {
	return RoleMetadata{
		Name:        "Default",
		Description: "Default role granted to all users in the organization.",
		IsDefault:   true,
		OrgId:       orgId,
	}
}

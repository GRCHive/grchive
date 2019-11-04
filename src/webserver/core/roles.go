package core

type AccessType struct {
	// View: Can see the thing being granted access to.
	CanView bool
	// Edit: Can change the thing being granted access to.
	CanEdit bool
	// Manage: Can add/delete the thing being granted access to.
	CanManage bool
}

type Permissions struct {
	// In the case of manage, allows giving/revoking a role from a user.
	OrgRoles AccessType
	// In this case, the editing the process flow includes adding/creating nodes.
	// Managing is merely for deleting/creating new process flows.
	ProcessFlows         AccessType
	Controls             AccessType
	ControlDocumentation AccessType
	Risks                AccessType
}

type Role struct {
	UserId           int64
	OrgToPermissions map[int32]Permissions
}

func CreateOwnerAccessType() AccessType {
	return AccessType{
		CanView:   true,
		CanEdit:   true,
		CanManage: true,
	}
}

func CreateAllAccessPermission() Permissions {
	return Permissions{
		OrgRoles:             CreateOwnerAccessType(),
		ProcessFlows:         CreateOwnerAccessType(),
		Controls:             CreateOwnerAccessType(),
		ControlDocumentation: CreateOwnerAccessType(),
		Risks:                CreateOwnerAccessType(),
	}
}

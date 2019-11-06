package core

type AccessType int

const (
	// 0b000
	AccessNone AccessType = 0
	// 0b001
	AccessView AccessType = 1
	// 0b010
	AccessEdit AccessType = 2
	// 0b100
	AccessManage AccessType = 4
)

type ResourceType int

const (
	// Managing is merely for deleting/creating new process flows.
	ResourceOrgRoles ResourceType = iota
	// In the case of manage, allows giving/revoking a role from a user.
	// In this case, the editing the process flow includes adding/creating nodes.
	ResourceProcessFlows
	ResourceControls
	ResourceControlDocumentation
	ResourceRisks
)

var AvailableResources []ResourceType = []ResourceType{
	ResourceOrgRoles,
	ResourceProcessFlows,
	ResourceControls,
	ResourceControlDocumentation,
	ResourceRisks,
}

type PermissionsMap struct {
	OrgRolesAccess             AccessType `db:"org_access"`
	ProcessFlowsAccess         AccessType `db:"flow_access"`
	ControlsAccess             AccessType `db:"control_access"`
	ControlDocumentationAccess AccessType `db:"doc_access"`
	RisksAccess                AccessType `db:"risk_access"`
}

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
	return AccessView | AccessEdit | AccessManage
}

func CreateAllAccessPermission() PermissionsMap {
	return PermissionsMap{
		OrgRolesAccess:             CreateOwnerAccessType(),
		ProcessFlowsAccess:         CreateOwnerAccessType(),
		ControlsAccess:             CreateOwnerAccessType(),
		ControlDocumentationAccess: CreateOwnerAccessType(),
		RisksAccess:                CreateOwnerAccessType(),
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

func (p PermissionsMap) GetAccessType(resource ResourceType) AccessType {
	switch resource {
	case ResourceOrgRoles:
		return p.OrgRolesAccess
	case ResourceProcessFlows:
		return p.ProcessFlowsAccess
	case ResourceControls:
		return p.ControlsAccess
	case ResourceControlDocumentation:
		return p.ControlDocumentationAccess
	case ResourceRisks:
		return p.RisksAccess
	}
	return AccessNone
}

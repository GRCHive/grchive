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
	ResourceOrgUsers ResourceType = iota
	// In the case of manage, allows giving/revoking a role from a user.
	ResourceOrgRoles
	// Managing is merely for deleting/creating new process flows.
	// In this case, the editing the process flow includes adding/creating nodes.
	ResourceProcessFlows
	ResourceControls
	ResourceControlDocumentation
	ResourceControlDocumentationMetadata
	ResourceRisks
	ResourceGeneralLedger
	ResourceSystems
	ResourceDatabases
	ResourceDbConnections
	ResourceDocRequests
	ResourceDeployments
	ResourceServers
	ResourceVendors
	ResourceDbSql
	ResourceDbSqlQuery
	ResourceDbSqlRequest
)

var AvailableResources []ResourceType = []ResourceType{
	ResourceOrgUsers,
	ResourceOrgRoles,
	ResourceProcessFlows,
	ResourceControls,
	ResourceControlDocumentation,
	ResourceControlDocumentationMetadata,
	ResourceRisks,
	ResourceGeneralLedger,
	ResourceSystems,
	ResourceDatabases,
	ResourceDbConnections,
	ResourceDocRequests,
	ResourceDeployments,
	ResourceServers,
	ResourceVendors,
	ResourceDbSql,
	ResourceDbSqlQuery,
	ResourceDbSqlRequest,
}

type PermissionsMap struct {
	OrgUsersAccess             AccessType `db:"users_access"`
	OrgRolesAccess             AccessType `db:"roles_access"`
	ProcessFlowsAccess         AccessType `db:"flow_access"`
	ControlsAccess             AccessType `db:"control_access"`
	ControlDocumentationAccess AccessType `db:"doc_access"`
	ControlDocMetadataAccess   AccessType `db:"doc_meta_access"`
	RisksAccess                AccessType `db:"risk_access"`
	GLAccess                   AccessType `db:"gl_access"`
	SystemAccess               AccessType `db:"system_access"`
	DbAccess                   AccessType `db:"db_access"`
	DbConnectionAccess         AccessType `db:"db_conn_access"`
	DocRequestAccess           AccessType `db:"doc_request_access"`
	DeploymentAccess           AccessType `db:"deployment_access"`
	ServerAccess               AccessType `db:"server_access"`
	VendorAccess               AccessType `db:"vendor_access"`
	DbSqlAccess                AccessType `db:"db_sql_access"`
	DbSqlQueryAccess           AccessType `db:"db_sql_query_access"`
	DbSqlRequestAccess         AccessType `db:"db_sql_requests_access"`
}

type RoleMetadata struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IsDefault   bool   `db:"is_default_role"`
	IsAdmin     bool   `db:"is_admin_role"`
	OrgId       int32  `db:"org_id"`
}

type Role struct {
	RoleMetadata `db:"role" json:"RoleMetadata"`
	Permissions  PermissionsMap `db:"permissions"`
	UserId       int64
}

func CreateOwnerAccessType() AccessType {
	return AccessView | AccessEdit | AccessManage
}

func CreateViewOnlyAccessPermission() PermissionsMap {
	return PermissionsMap{
		OrgUsersAccess:             AccessView,
		OrgRolesAccess:             AccessView,
		ProcessFlowsAccess:         AccessView,
		ControlsAccess:             AccessView,
		ControlDocumentationAccess: AccessView,
		ControlDocMetadataAccess:   AccessView,
		RisksAccess:                AccessView,
		GLAccess:                   AccessView,
		SystemAccess:               AccessView,
		DbAccess:                   AccessView,
		DbConnectionAccess:         AccessNone,
		DocRequestAccess:           AccessView,
		DeploymentAccess:           AccessView,
		ServerAccess:               AccessView,
		VendorAccess:               AccessView,
		DbSqlAccess:                AccessView,
		DbSqlQueryAccess:           AccessView,
		DbSqlRequestAccess:         AccessView,
	}
}

func CreateAllAccessPermission() PermissionsMap {
	return PermissionsMap{
		OrgUsersAccess:             CreateOwnerAccessType(),
		OrgRolesAccess:             CreateOwnerAccessType(),
		ProcessFlowsAccess:         CreateOwnerAccessType(),
		ControlsAccess:             CreateOwnerAccessType(),
		ControlDocumentationAccess: CreateOwnerAccessType(),
		ControlDocMetadataAccess:   CreateOwnerAccessType(),
		RisksAccess:                CreateOwnerAccessType(),
		GLAccess:                   CreateOwnerAccessType(),
		SystemAccess:               CreateOwnerAccessType(),
		DbAccess:                   CreateOwnerAccessType(),
		DbConnectionAccess:         CreateOwnerAccessType(),
		DocRequestAccess:           CreateOwnerAccessType(),
		DeploymentAccess:           CreateOwnerAccessType(),
		ServerAccess:               CreateOwnerAccessType(),
		VendorAccess:               CreateOwnerAccessType(),
		DbSqlAccess:                CreateOwnerAccessType(),
		DbSqlQueryAccess:           CreateOwnerAccessType(),
		DbSqlRequestAccess:         CreateOwnerAccessType(),
	}
}

var ServerRole = &Role{
	RoleMetadata: RoleMetadata{
		Id:          -1,
		Name:        "Server Role",
		Description: "All access Server role",
	},
	Permissions: CreateAllAccessPermission(),
}

func CreateDefaultRoleMetadata(orgId int32) RoleMetadata {
	return RoleMetadata{
		Name:        "Default",
		Description: "Default role granted to all users in the organization.",
		IsDefault:   true,
		IsAdmin:     false,
		OrgId:       orgId,
	}
}

func CreateAdminRoleMetadata(orgId int32) RoleMetadata {
	return RoleMetadata{
		Name:        "Admin",
		Description: "Admin role for the organization.",
		IsDefault:   false,
		IsAdmin:     true,
		OrgId:       orgId,
	}
}

func (p PermissionsMap) GetAccessType(resource ResourceType) AccessType {
	switch resource {
	case ResourceOrgUsers:
		return p.OrgUsersAccess
	case ResourceOrgRoles:
		return p.OrgRolesAccess
	case ResourceProcessFlows:
		return p.ProcessFlowsAccess
	case ResourceControls:
		return p.ControlsAccess
	case ResourceControlDocumentation:
		return p.ControlDocumentationAccess
	case ResourceControlDocumentationMetadata:
		return p.ControlDocMetadataAccess
	case ResourceRisks:
		return p.RisksAccess
	case ResourceGeneralLedger:
		return p.GLAccess
	case ResourceSystems:
		return p.SystemAccess
	case ResourceDatabases:
		return p.DbAccess
	case ResourceDbConnections:
		return p.DbConnectionAccess
	case ResourceDocRequests:
		return p.DocRequestAccess
	case ResourceDeployments:
		return p.DeploymentAccess
	case ResourceServers:
		return p.ServerAccess
	case ResourceVendors:
		return p.VendorAccess
	case ResourceDbSql:
		return p.DbSqlAccess
	case ResourceDbSqlQuery:
		return p.DbSqlQueryAccess
	case ResourceDbSqlRequest:
		return p.DbSqlRequestAccess
	}
	return AccessNone
}

func (p *PermissionsMap) SetAccessType(resource ResourceType, access AccessType) {
	switch resource {
	case ResourceOrgUsers:
		p.OrgUsersAccess = access
		break
	case ResourceOrgRoles:
		p.OrgRolesAccess = access
		break
	case ResourceProcessFlows:
		p.ProcessFlowsAccess = access
		break
	case ResourceControls:
		p.ControlsAccess = access
		break
	case ResourceControlDocumentation:
		p.ControlDocumentationAccess = access
		break
	case ResourceControlDocumentationMetadata:
		p.ControlDocMetadataAccess = access
		break
	case ResourceRisks:
		p.RisksAccess = access
		break
	case ResourceGeneralLedger:
		p.GLAccess = access
		break
	case ResourceSystems:
		p.SystemAccess = access
		break
	case ResourceDatabases:
		p.DbAccess = access
		break
	case ResourceDbConnections:
		p.DbConnectionAccess = access
		break
	case ResourceDocRequests:
		p.DocRequestAccess = access
		break
	case ResourceDeployments:
		p.DeploymentAccess = access
		break
	case ResourceServers:
		p.ServerAccess = access
		break
	case ResourceVendors:
		p.VendorAccess = access
		break
	case ResourceDbSql:
		p.DbSqlAccess = access
		break
	case ResourceDbSqlQuery:
		p.DbSqlQueryAccess = access
	case ResourceDbSqlRequest:
		p.DbSqlRequestAccess = access
	}
}

func (p PermissionsMap) HasAccess(resource ResourceType, access AccessType) bool {
	return (p.GetAccessType(resource) & access) != 0
}

func MaxAccessType(a AccessType, b AccessType) AccessType {
	if a < b {
		return b
	}
	return a
}

func MinAccessType(a AccessType, b AccessType) AccessType {
	if a > b {
		return b
	}
	return a
}

func (r *Role) SetMax(access AccessType) {
	for _, resource := range AvailableResources {
		r.Permissions.SetAccessType(
			resource,
			MinAccessType(access, r.Permissions.GetAccessType(resource)))
	}
}

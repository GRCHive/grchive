export interface RoleMetadata {
    Id: number
    Name: string
    Description: string
    IsDefault: boolean
    IsAdmin: boolean
    OrgId: number
}

export enum AccessType {
    NoAccess = 0,
    View = 1,
    Edit = 2,
    Manage = 4
}

export interface Permissions {
	OrgUsersAccess            : AccessType 
	OrgRolesAccess            : AccessType 
	ProcessFlowsAccess        : AccessType 
	ControlsAccess            : AccessType 
	ControlDocumentationAccess: AccessType 
	ControlDocMetadataAccess  : AccessType 
	RisksAccess               : AccessType 
    GLAccess                  : AccessType
    SystemAccess              : AccessType
    DbAccess                  : AccessType
    DbConnectionAccess        : AccessType
    DocRequestAccess          : AccessType
    DeploymentAccess          : AccessType
    ServerAccess              : AccessType
    VendorAccess              : AccessType
    DbSqlAccess               : AccessType
    DbSqlQueryAccess          : AccessType
    DbSqlRequestAccess        : AccessType
	ClientDataAccess          : AccessType
	ManagedCodeAccess         : AccessType
    ClientScriptsAccess       : AccessType
    ScriptRunAccess           : AccessType
    BuildLogAccess            : AccessType
    ShellScriptAccess         : AccessType
    ShellScriptRunsAccess     : AccessType
    IntegrationAccess         : AccessType
    SapErpAccess              : AccessType
}

export interface FullRole {
    RoleMetadata: RoleMetadata
    Permissions: Permissions
}

export function emptyPermissions() : Permissions {
    return {
        OrgUsersAccess: AccessType.NoAccess,
        OrgRolesAccess: AccessType.NoAccess,
        ProcessFlowsAccess: AccessType.NoAccess,
        ControlsAccess: AccessType.NoAccess,
        ControlDocumentationAccess: AccessType.NoAccess,
        ControlDocMetadataAccess: AccessType.NoAccess,
        RisksAccess: AccessType.NoAccess,
        GLAccess: AccessType.NoAccess,
        SystemAccess: AccessType.NoAccess,
        DbAccess: AccessType.NoAccess,
        DbConnectionAccess: AccessType.NoAccess,
        DocRequestAccess: AccessType.NoAccess,
        DeploymentAccess: AccessType.NoAccess,
        ServerAccess: AccessType.NoAccess,
        VendorAccess: AccessType.NoAccess,
        DbSqlAccess: AccessType.NoAccess,
        DbSqlQueryAccess: AccessType.NoAccess,
        DbSqlRequestAccess: AccessType.NoAccess,
        ClientDataAccess: AccessType.NoAccess,
        ManagedCodeAccess: AccessType.NoAccess,
        ClientScriptsAccess: AccessType.NoAccess,
        ScriptRunAccess: AccessType.NoAccess,
        BuildLogAccess: AccessType.NoAccess,
        ShellScriptAccess: AccessType.NoAccess,
        ShellScriptRunsAccess: AccessType.NoAccess,
        IntegrationAccess: AccessType.NoAccess,
        SapErpAccess: AccessType.NoAccess,
    }
}

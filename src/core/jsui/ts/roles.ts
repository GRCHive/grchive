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
}

export interface FullRole {
    RoleMetadata: RoleMetadata
    Permissions: Permissions
}

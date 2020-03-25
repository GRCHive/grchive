package grchive.core.data.types.grchive

/**
 * Contains information about the role; note that this class does not
 * contain information about what permissions are stored.
 *
 * @property id Unique identifier.
 * @property name Human readable name.
 * @property description An optional description.
 * @property isDefault If true, this is the default role given to new users to the organization.
 * @property isAdmin If true, this is an admin role and should have all permissions.
 * @property orgId The organization this role belongs to.
 */
data class Role (
    val id : Long,
    val name : String,
    val description : String,
    val isDefault: Boolean,
    val isAdmin: Boolean,
    val orgId : Int
)

/**
 * Contains both the [Role] and its corresponding [RolePermissions].
 */
data class FullRole (
    val role : Role,
    val permissions : RolePermissions
)

package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName
import org.jdbi.v3.core.mapper.Nested

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
    @ColumnName("id") val id : Long,
    @ColumnName("name") val name : String,
    @ColumnName("description") val description : String,
    @ColumnName("is_default_role") val isDefault: Boolean,
    @ColumnName("is_admin_role") val isAdmin: Boolean,
    @ColumnName("org_id") val orgId : Int
)

/**
 * Contains both the [Role] and its corresponding [RolePermissions].
 */
data class FullRole (
    @Nested("role") val role : Role,
    @Nested("permissions") val permissions : RolePermissions
)

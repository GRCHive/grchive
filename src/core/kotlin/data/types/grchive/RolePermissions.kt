package grchive.core.data.types.grchive

/**
 * Unionable permissions to grant to a role for accessing a resource.
 */
enum class AccessType(val bit : Int) {
    None(0b000),
    View(0b001),
    Edit(0b010),
    Manage(0b100)
}

/**
 * Permissions granted to a [Role] specified as a union
 * of [AccessType].
 *
 * @property orgUsersAccess
 * @property orgRolesAccess
 * @property processFlowsAccess
 * @property controlsAccess
 * @property controlDocumentationAccess
 * @property controlDocMetadataAccess
 * @property risksAccess
 * @property gLAccess
 * @property systemAccess
 * @property dbAccess
 * @property dbConnectionAccess
 * @property docRequestAccess
 * @property deploymentAccess
 * @property serverAccess
 * @property vendorAccess
 * @property dbSqlAccess
 * @property dbSqlQueryAccess
 * @property dbSqlRequestAccess
 */
data class RolePermissions (
    val orgUsersAccess             : Int,
    val orgRolesAccess             : Int,
    val processFlowsAccess         : Int,
    val controlsAccess             : Int,
    val controlDocumentationAccess : Int,
    val controlDocMetadataAccess   : Int,
    val risksAccess                : Int,
    val gLAccess                   : Int,
    val systemAccess               : Int,
    val dbAccess                   : Int,
    val dbConnectionAccess         : Int,
    val docRequestAccess           : Int,
    val deploymentAccess           : Int,
    val serverAccess               : Int,
    val vendorAccess               : Int,
    val dbSqlAccess                : Int,
    val dbSqlQueryAccess           : Int,
    val dbSqlRequestAccess         : Int 
)

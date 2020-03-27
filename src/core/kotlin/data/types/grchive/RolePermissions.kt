package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * Unionable permissions to grant to a role for accessing a resource.
 */
enum class AccessType(val bit : Int) {
    None(0b000),
    View(0b001),
    Edit(0b010),
    Manage(0b100),
    All(0b111)
}

/**
 * The exception that gets thrown whenever a resource is accessed without the role having proper permissions.
 *
 * @property res The resource requested.
 * @property access The access requested.
 */
class ResourcePermissionDeniedException(val res : Resources, val access : AccessType) : Exception("Permission denied for resource access.")

/**
 * Unions multiple access types together to create a unioned access type.
 *
 * @param v One or more [AccessType] enums
 * @return An Int that contains the unioned bit values.
 */
fun unionAccessType(vararg v : AccessType) : Int {
    var ret : Int = 0
    v.forEach {
        ret = ret or it.bit
    }
    return ret
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
 * @property clientDataAccess
 */
data class RolePermissions (
    @field:GrchiveResource(Resources.ResourceOrgUsers) @ColumnName("users_access")
    val orgUsersAccess : Int,
    @field:GrchiveResource(Resources.ResourceOrgRoles) @ColumnName("roles_access")
    val orgRolesAccess : Int,
    @field:GrchiveResource(Resources.ResourceProcessFlows) @ColumnName("flow_access")
    val processFlowsAccess : Int,
    @field:GrchiveResource(Resources.ResourceControls) @ColumnName("control_access")
    val controlsAccess : Int,
    @field:GrchiveResource(Resources.ResourceControlDocumentation) @ColumnName("doc_access")
    val controlDocumentationAccess : Int,
    @field:GrchiveResource(Resources.ResourceControlDocumentationMetadata) @ColumnName("doc_meta_access")
    val controlDocMetadataAccess : Int,
    @field:GrchiveResource(Resources.ResourceRisks) @ColumnName("risk_access")
    val risksAccess: Int,
    @field:GrchiveResource(Resources.ResourceGeneralLedger) @ColumnName("gl_access")
    val glAccess : Int,
    @field:GrchiveResource(Resources.ResourceSystems) @ColumnName("system_access")
    val systemAccess : Int,
    @field:GrchiveResource(Resources.ResourceDatabases) @ColumnName("db_access")
    val dbAccess : Int,
    @field:GrchiveResource(Resources.ResourceDbConnections) @ColumnName("db_conn_access")
    val dbConnectionAccess : Int,
    @field:GrchiveResource(Resources.ResourceDocRequests) @ColumnName("doc_request_access")
    val docRequestAccess : Int,
    @field:GrchiveResource(Resources.ResourceDeployments) @ColumnName("deployment_access")
    val deploymentAccess : Int,
    @field:GrchiveResource(Resources.ResourceServers) @ColumnName("server_access")
    val serverAccess : Int,
    @field:GrchiveResource(Resources.ResourceVendors) @ColumnName("vendor_access")
    val vendorAccess : Int,
    @field:GrchiveResource(Resources.ResourceDbSql) @ColumnName("db_sql_access")
    val dbSqlAccess: Int,
    @field:GrchiveResource(Resources.ResourceDbSqlQuery) @ColumnName("db_sql_query_access")
    val dbSqlQueryAccess : Int,
    @field:GrchiveResource(Resources.ResourceDbSqlRequest) @ColumnName("db_sql_requests_access")
    val dbSqlRequestAccess : Int
    @field:GrchiveResource(Resources.ResourceClientData) @ColumnName("client_data_access")
    val clientDataAccess : Int
)

fun emptyRolePermissions() : RolePermissions {
    return RolePermissions(
        AccessType.None.bit /* orgUsersAccess */,
        AccessType.None.bit /* orgRolesAccess */,
        AccessType.None.bit /* processFlowsAccess */,
        AccessType.None.bit /* controlsAccess */,
        AccessType.None.bit /* controlDocumentationAccess */,
        AccessType.None.bit /* controlDocMetadataAccess */,
        AccessType.None.bit /* risksAccess*/,
        AccessType.None.bit /* glAccess */,
        AccessType.None.bit /* systemAccess */,
        AccessType.None.bit /* dbAccess */,
        AccessType.None.bit /* dbConnectionAccess */,
        AccessType.None.bit /* docRequestAccess */,
        AccessType.None.bit /* deploymentAccess */,
        AccessType.None.bit /* serverAccess */,
        AccessType.None.bit /* vendorAccess */,
        AccessType.None.bit /* dbSqlAccess*/,
        AccessType.None.bit /* dbSqlQueryAccess */,
        AccessType.None.bit /* dbSqlRequestAccess */
        AccessType.None.bit /* clientDataAccess */,
    )
}

fun fullRolePermissions() : RolePermissions {
    return RolePermissions(
        AccessType.All.bit /* orgUsersAccess */,
        AccessType.All.bit /* orgRolesAccess */,
        AccessType.All.bit /* processFlowsAccess */,
        AccessType.All.bit /* controlsAccess */,
        AccessType.All.bit /* controlDocumentationAccess */,
        AccessType.All.bit /* controlDocMetadataAccess */,
        AccessType.All.bit /* risksAccess*/,
        AccessType.All.bit /* glAccess */,
        AccessType.All.bit /* systemAccess */,
        AccessType.All.bit /* dbAccess */,
        AccessType.All.bit /* dbConnectionAccess */,
        AccessType.All.bit /* docRequestAccess */,
        AccessType.All.bit /* deploymentAccess */,
        AccessType.All.bit /* serverAccess */,
        AccessType.All.bit /* vendorAccess */,
        AccessType.All.bit /* dbSqlAccess*/,
        AccessType.All.bit /* dbSqlQueryAccess */,
        AccessType.All.bit /* dbSqlRequestAccess */,
        AccessType.All.bit /* clientDataAccess */,
    )
}

/**
 * Returns the unioned access type for the given resource in the permission map.
 *
 * @param p The [RolePermissions] map to retrieve from.
 * @param r The [Resources] to access.
 * @return A unioned [AccessType].
 */
fun getRolePermissionForResource(p : RolePermissions, r : Resources) : Int {
    return when (r) {
        Resources.ResourceOrgUsers -> p.orgUsersAccess
        Resources.ResourceOrgRoles -> p.orgRolesAccess
        Resources.ResourceProcessFlows -> p.processFlowsAccess
        Resources.ResourceControls -> p.controlsAccess
        Resources.ResourceControlDocumentation -> p.controlDocumentationAccess
        Resources.ResourceControlDocumentationMetadata -> p.controlDocMetadataAccess
        Resources.ResourceRisks -> p.risksAccess
        Resources.ResourceGeneralLedger -> p.glAccess
        Resources.ResourceSystems -> p.systemAccess
        Resources.ResourceDatabases -> p.dbAccess
        Resources.ResourceDbConnections -> p.dbConnectionAccess
        Resources.ResourceDocRequests -> p.docRequestAccess
        Resources.ResourceDeployments -> p.deploymentAccess
        Resources.ResourceServers -> p.serverAccess
        Resources.ResourceVendors -> p.vendorAccess
        Resources.ResourceDbSql -> p.dbSqlAccess
        Resources.ResourceDbSqlQuery -> p.dbSqlQueryAccess
        Resources.ResourceDbSqlRequest -> p.dbSqlRequestAccess
        Resources.ResourceClientData -> p.clientDataAccess
    }
}

/**
 * Utility function to check if a role has a certain type of access to a resource.
 *
 * @param p The [RolePermissions] object to check.
 * @param r The [Resources] to check access for.
 * @param a The type of desired [AccessType].
 * @return True if the role has the requested access.
 */
fun roleHasAccess(p : RolePermissions, r : Resources, a : AccessType) : Boolean {
    return (getRolePermissionForResource(p, r) and a.bit) != 0
}

/**
 * Check if the role has access, if not, throw an exception of type [ResourcePermissionDeniedException].
 *
 * @param p The [RolePermissions] object to check.
 * @param r The [Resources] to check access for.
 * @param a The type of desired [AccessType].
 */
fun roleMustHaveAccess(p : RolePermissions, r : Resources, a : AccessType) {
    if (!roleHasAccess(p, r, a)) {
        throw ResourcePermissionDeniedException(r, a)
    }
}

package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

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
 */
data class RolePermissions (
    @field:GrchiveResource(Resources.ResourceOrgUsers) @ColumnName("users_access")
    val orgUsersAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceOrgRoles) @ColumnName("roles_access")
    val orgRolesAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceProcessFlows) @ColumnName("flow_access")
    val processFlowsAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceControls) @ColumnName("control_access")
    val controlsAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceControlDocumentation) @ColumnName("doc_access")
    val controlDocumentationAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceControlDocumentationMetadata) @ColumnName("doc_meta_access")
    val controlDocMetadataAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceRisks) @ColumnName("risk_access")
    val risksAccess: Int = 0,
    @field:GrchiveResource(Resources.ResourceGeneralLedger) @ColumnName("gl_access")
    val glAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceSystems) @ColumnName("system_access")
    val systemAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDatabases) @ColumnName("db_access")
    val dbAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDbConnections) @ColumnName("db_conn_access")
    val dbConnectionAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDocRequests) @ColumnName("doc_request_access")
    val docRequestAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDeployments) @ColumnName("deployment_access")
    val deploymentAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceServers) @ColumnName("server_access")
    val serverAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceVendors) @ColumnName("vendor_access")
    val vendorAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDbSql) @ColumnName("db_sql_access")
    val dbSqlAccess: Int = 0,
    @field:GrchiveResource(Resources.ResourceDbSqlQuery) @ColumnName("db_sql_query_access")
    val dbSqlQueryAccess : Int = 0,
    @field:GrchiveResource(Resources.ResourceDbSqlRequest) @ColumnName("db_sql_requests_access")
    val dbSqlRequestAccess : Int = 0
)

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
    }
}

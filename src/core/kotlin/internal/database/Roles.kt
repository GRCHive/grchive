package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import java.util.Optional
import java.lang.StringBuilder

import kotlin.collections.mapOf

import grchive.core.data.types.grchive.FullRole
import grchive.core.data.types.grchive.Resources

internal val resourceToDatabaseMap = mapOf(
	Resources.ResourceOrgUsers to                     "resource_organization_users_access",
	Resources.ResourceOrgRoles to                     "resource_organization_roles_access",
	Resources.ResourceProcessFlows to                 "resource_process_flows_access",
	Resources.ResourceControls to                     "resource_controls_access",
	Resources.ResourceControlDocumentation to         "resource_control_documentation_access",
	Resources.ResourceControlDocumentationMetadata to "resource_control_documentation_metadata_access",
	Resources.ResourceRisks to                        "resource_risks_access",
	Resources.ResourceGeneralLedger to                "resource_gl_access",
	Resources.ResourceSystems to                      "resource_systems_access",
	Resources.ResourceDatabases to                    "resource_database_access",
	Resources.ResourceDbConnections to                "resource_db_conn_access",
	Resources.ResourceDocRequests to                  "resource_doc_request_access",
	Resources.ResourceDeployments to                  "resource_deployment_access",
	Resources.ResourceServers to                      "resource_server_access",
	Resources.ResourceVendors to                      "resource_vendor_access",
	Resources.ResourceDbSql to                        "resource_db_sql_access",
	Resources.ResourceDbSqlQuery to                   "resource_db_sql_query_access",
	Resources.ResourceDbSqlRequest to                 "resource_db_sql_requests_access",
    Resources.ResourceClientData to                   "resource_client_data_access",
    Resources.ResourceManagedCode to                  "resource_managed_code_access",
    Resources.ResourceClientScripts to                "resource_client_scripts_access"
)

internal val resourceToColumnName = mapOf(
	Resources.ResourceOrgUsers to                     "permissions_users_access",
	Resources.ResourceOrgRoles to                     "permissions_roles_access",
	Resources.ResourceProcessFlows to                 "permissions_flow_access",
	Resources.ResourceControls to                     "permissions_control_access",
	Resources.ResourceControlDocumentation to         "permissions_doc_access",
	Resources.ResourceControlDocumentationMetadata to "permissions_doc_meta_access",
	Resources.ResourceRisks to                        "permissions_risk_access",
	Resources.ResourceGeneralLedger to                "permissions_gl_access",
	Resources.ResourceSystems to                      "permissions_system_access",
	Resources.ResourceDatabases to                    "permissions_db_access",
	Resources.ResourceDbConnections to                "permissions_db_conn_access",
	Resources.ResourceDocRequests to                  "permissions_doc_request_access",
	Resources.ResourceDeployments to                  "permissions_deployment_access",
	Resources.ResourceServers to                      "permissions_server_access",
	Resources.ResourceVendors to                      "permissions_vendor_access",
	Resources.ResourceDbSql to                        "permissions_db_sql_access",
	Resources.ResourceDbSqlQuery to                   "permissions_db_sql_query_access",
	Resources.ResourceDbSqlRequest to                 "permissions_db_sql_requests_access",
    Resources.ResourceClientData to                   "permissions_client_data_access",
    Resources.ResourceManagedCode to                  "permissions_managed_code_access",
    Resources.ResourceClientScripts to                "permissions_client_scripts_access"
)

private fun createResourceSelect() : String {
    val sb = StringBuilder()

    enumValues<Resources>().forEach {
        sb.append(", ${resourceToDatabaseMap[it]}.access_type AS \"${resourceToColumnName[it]}\"\n")
    }

    return sb.toString()
}

private fun createResourceJoin() : String {
    val sb = StringBuilder()

    enumValues<Resources>().forEach {
        val tb = resourceToDatabaseMap[it]
        sb.append("LEFT JOIN ${tb} ON role.id = ${tb}.role_id\n")
    }

    return sb.toString()
}

/**
 * Finds the [FullRole] for the given user in the organization.
 *
 * @param hd A JDBI handle.
 * @param userId The User ID with an assigned role.
 * @param orgId The org ID that the user belongs to.
 * @return A [FullRole] if found. Null if not.
 */
fun getUserRoleForOrg(hd : Handle, userId : Long, orgId : Int) : FullRole? {
    val res : Optional<FullRole> = hd.select("""
		SELECT
            role.id AS "role_id",
            role.name AS "role_name",
            role.description AS "role_description",
            role.is_default_role AS "role_is_default_role",
            role.is_admin_role AS "role_is_admin_role",
            role.org_id AS "role_org_id"
            ${createResourceSelect()}
		FROM organization_available_roles AS role
        ${createResourceJoin()}
        INNER JOIN user_roles AS ur
            ON ur.role_id = role.id
        WHERE ur.user_id = ? AND ur.org_id = ?
    """, userId, orgId).mapTo(FullRole::class.java).findOne()

    return if (res.isPresent()) res.get() else null
}

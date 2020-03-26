package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.FullRole
import grchive.core.data.types.grchive.ProcessFlowMetadata
import grchive.core.data.types.grchive.roleMustHaveAccess
import grchive.core.data.types.grchive.Resources
import grchive.core.data.types.grchive.AccessType
import grchive.core.data.filters.Filter

/**
 * Returns all process flow metadata that matches the given filters for the organization.
 *
 * @param hd A JDBI handle to the database.
 * @param orgId The unique ID for the organization.
 * @param role The [FullRole] that is being used to access the resource.
 * @param filters A map of [Filter]s to determine which resources to return. The map key is the column name.
 * @return A list of [ProcessFlowMetadata].
 */
internal fun getAllProcessFlowMetadata(
    hd : Handle,
    orgId : Int,
    role : FullRole,
    filters : Map<String, Filter>
) : List<ProcessFlowMetadata> {
    roleMustHaveAccess(role.permissions, Resources.ResourceProcessFlows, AccessType.View)

    val (condition, args) = createSqlConditionsFromFilters(filters, prefix="AND", columnPrefix="flow.")
    return hd.select("""
        SELECT flow.*
        FROM process_flows AS flow
        WHERE flow.org_id = ?
        ${condition}
    """, orgId, *args.toArray()).mapTo(ProcessFlowMetadata::class.java).list()
}

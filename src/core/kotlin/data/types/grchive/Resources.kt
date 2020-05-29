package grchive.core.data.types.grchive

enum class Resources {
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
    ResourceClientData,
    ResourceManagedCode,
    ResourceClientScripts,
    ResourceScriptRun,
    ResourceBuildLog,
    ResourceShell
}

/**
 * Allow for annotating something with the corresponding resource.
 */
@MustBeDocumented
annotation class GrchiveResource(val res : Resources)

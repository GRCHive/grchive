package grchive.core.data.sources

import grchive.core.data.types.grchive.ClientDataSourceLink
import grchive.core.data.types.grchive.SupportedDataSources
import grchive.core.data.types.grchive.ApiKey
import grchive.core.internal.Config

internal fun createGrchiveDataSource(cfg : Config, userId : Long, orgId : Int) : GrchiveDataSource {
    return GrchiveDataSource(cfg, userId, orgId)
}

internal fun makeDataSourceFromClientDataSourceLink(
    link : ClientDataSourceLink,
    cfg : Config,
    userId : Long,
    orgId : Int)
: RawDataSource {
    return when (link.sourceId) {
        SupportedDataSources.kGrchive.id -> createGrchiveDataSource(cfg, userId, orgId)
        else -> {
            throw Exception("Unsupported data source.")
        }
    }
}

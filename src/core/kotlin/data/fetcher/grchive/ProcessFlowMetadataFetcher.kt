package grchive.core.data.fetcher.grchive

import org.jdbi.v3.core.kotlin.*

import grchive.core.data.fetcher.DataFetcher
import grchive.core.data.filters.Filter
import grchive.core.data.sources.GrchiveDataSource
import grchive.core.data.types.grchive.ProcessFlowMetadata

import grchive.core.internal.database.getAllProcessFlowMetadata

final class ProcessFlowMetadataFetcher : DataFetcher<ProcessFlowMetadata, GrchiveDataSource> {
    override fun fetch(source : GrchiveDataSource, filters : Map<String, Filter>) : List<ProcessFlowMetadata> {
        return source.jdbi.withHandleUnchecked {
            getAllProcessFlowMetadata(it, source.orgId, source.activeRole, filters)
        }
    }
}
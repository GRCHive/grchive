package grchive.core.data.fetcher.grchive

import org.jdbi.v3.core.kotlin.*

import grchive.core.data.fetcher.DataFetcher
import grchive.core.data.filters.Filter
import grchive.core.data.sources.GrchiveDataSource
import grchive.core.data.track.TrackedData
import grchive.core.data.track.TrackedSource
import grchive.core.data.track.TrackedSourceLogger
import grchive.core.data.types.grchive.User

import grchive.core.internal.database.getUsersInOrganizationId

final class UsersFetcher : DataFetcher<User, GrchiveDataSource> {
    override fun fetch(source : GrchiveDataSource, filters : Map<String, Filter>) : List<TrackedData<User>> {
        val trackSrc = TrackedSource(source.clientDataId)
        val logger = TrackedSourceLogger(trackSrc)

        return source.db.jdbi.withHandleUnchecked {
            it.setSqlLogger(logger)
            getUsersInOrganizationId(it, source.orgId).map {
                TrackedData<User>(it, trackSrc)
            }
        }
    }
}

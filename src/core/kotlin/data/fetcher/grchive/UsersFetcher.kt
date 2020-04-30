package grchive.core.data.fetcher.grchive

import org.jdbi.v3.core.kotlin.*

import grchive.core.data.fetcher.DataFetcher
import grchive.core.data.filters.Filter
import grchive.core.data.sources.GrchiveDataSource
import grchive.core.data.types.grchive.User

import grchive.core.internal.database.getUsersInOrganizationId

final class UsersFetcher : DataFetcher<User, GrchiveDataSource> {
    override fun fetch(source : GrchiveDataSource, filters : Map<String, Filter>) : List<User> {
        return source.db.jdbi.withHandleUnchecked {
            getUsersInOrganizationId(it, source.orgId)
        }
    }
}

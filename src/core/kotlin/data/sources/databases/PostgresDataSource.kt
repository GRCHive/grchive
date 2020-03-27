package grchive.core.data.sources.databases

import javax.sql.DataSource

import org.jdbi.v3.postgres.PostgresPlugin

import grchive.core.data.sources.DatabaseDataSource
import grchive.core.utility.database.PostgresReadOnlyHandler

/** 
 * An abstraction over a read-only connection to a PostgreSQL database.
 */
open class PostgresDataSource(ds : DataSource) : DatabaseDataSource(ds, PostgresReadOnlyHandler()) {
    init {
        jdbi.installPlugin(PostgresPlugin())
    }
}

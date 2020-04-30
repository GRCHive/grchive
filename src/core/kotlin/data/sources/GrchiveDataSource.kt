package grchive.core.data.sources

import grchive.core.data.sources.databases.PostgresDataSource
import grchive.core.data.types.grchive.*

import grchive.core.internal.Config
import grchive.core.internal.database.*
import org.jdbi.v3.core.Handle
import org.jdbi.v3.core.kotlin.*

/** 
 * An abstraction over the connection to the data (databases, documentation, etc.) stored on GRCHive.
 *
 * @property cfg The GRCHive configuration to access the data sources.
 * @property apiKey The API key used to authorize actions on GRCHive data.
 * @property orgId The unique ID of the organization to access data for.
 */
class GrchiveDataSource internal constructor(
    internal val cfg : Config,
    val userId : Long,
    val orgId : Int) : RawDataSource {

    val activeRole : FullRole
    internal val db = PostgresDataSource(createGrchiveHikariDataSource(cfg.database, true))

    init {
        setupGrchiveJdbi(db.jdbi)

        activeRole = db.withHandle {
            getUserRoleForOrg(it, userId, orgId)!!
        }
    }
}

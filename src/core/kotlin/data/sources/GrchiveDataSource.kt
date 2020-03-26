package grchive.core.data.sources

import grchive.core.data.types.grchive.*

import grchive.core.internal.Config
import grchive.core.internal.database.*
import grchive.core.utility.database.PostgresReadOnlyHandler
import org.jdbi.v3.core.Handle
import org.jdbi.v3.core.kotlin.*

/** 
 * An abstraction over the connection to the data (databases, documentation, etc.) stored on GRCHive.
 *
 * @property cfg The GRCHive configuration to access the data sources.
 * @property apiKey The API key used to authorize actions on GRCHive data.
 * @property orgId The unique ID of the organization to access data for.
 */
class GrchiveDataSource(val cfg : Config, val apiKey : String, val orgId : Int) 
    : DatabaseDataSource(createGrchiveHikariDataSource(cfg.database, true), PostgresReadOnlyHandler()) {
    val activeRole : FullRole

    init {
        setupGrchiveJdbi(jdbi)

        activeRole = withHandle {
            val key : ApiKey = getApiKeyFromRawKey(it, apiKey)!!
            getRoleAttachedToApiKey(it, key.id, orgId)!!
        }
    }
}

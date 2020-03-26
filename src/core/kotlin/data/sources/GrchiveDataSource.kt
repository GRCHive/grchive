package grchive.core.data.sources

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
 */
class GrchiveDataSource(val cfg : Config, val apiKey : String) 
    : DatabaseDataSource(createGrchiveHikariDataSource(cfg.database)) {
    val activeRole : FullRole

    init {
        setupGrchiveJdbi(jdbi)

        activeRole = jdbi.withHandleUnchecked {
            val key : ApiKey = getApiKeyFromRawKey(it, apiKey)!!
            getRoleAttachedToApiKey(it, key.id)!!
        }
    }
}

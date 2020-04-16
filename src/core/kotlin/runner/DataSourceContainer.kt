package grchive.core.runner

import org.jdbi.v3.core.Handle

import grchive.core.data.sources.RawDataSource
import grchive.core.data.sources.makeDataSourceFromClientDataSourceLink

import grchive.core.data.types.grchive.ApiKey
import grchive.core.internal.Config
import grchive.core.internal.database.getClientDataFromId
import grchive.core.internal.database.getClientDataSourceLinkFromDataId

/**
 * A container that stores client data id, data source pairs (this is not Java's data source, but rather our own).
 */
class DataSourceContainer {
    // Meh on making this public but making it protected causes getSource to compile with a warning so here we are...
    val sources : MutableMap<String, RawDataSource> = mutableMapOf<String, RawDataSource>()

    internal fun addSource(k : String, s : RawDataSource) {
        sources.put(k, s)
    }

    fun keys() : Set<String> {
        return sources.keys
    }

    inline fun <reified T : RawDataSource> getSource(k : String) : T? {
        val ds = sources.get(k)
        if (ds == null) {
            return null
        }

        if (ds is T) {
            return ds
        }

        return null
    }
}

/**
 * Creates a [DataSourceContainer].
 *
 * @return A [DataSourceContainer] that holds a properly subtyped [RawDataSource] for every client data.
 * @param handle A JDBI handle to connect to the GRCHive database.
 * @param meta A [Metadata] object that holds what client data we need to load.
 */
internal fun loadDataSourceContainer(handle : Handle, meta : Metadata, cfg : Config, userId : Long, orgId : Int) : DataSourceContainer {
    val container = DataSourceContainer()

    meta.clientDataId.forEach {
        val clientData = getClientDataFromId(handle, it)
        val source = getClientDataSourceLinkFromDataId(handle, it)
        container.addSource(clientData.name, makeDataSourceFromClientDataSourceLink(source, cfg, userId, orgId))
    }

    return container
}

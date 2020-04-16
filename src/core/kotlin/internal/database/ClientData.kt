package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.ClientData
import grchive.core.data.types.grchive.ClientDataSourceLink

internal fun getClientDataFromId(
    hd : Handle,
    id : Long
) : ClientData {
    return hd.select("""
        SELECT *
        FROM client_data
        WHERE id = ?
    """, id).mapTo(ClientData::class.java).one()
}

internal fun getClientDataSourceLinkFromDataId(
    hd : Handle,
    dataId : Long
) : ClientDataSourceLink {
    return hd.select("""
        SELECT *
        FROM client_data_source_link
        WHERE data_id = ?
    """, dataId).mapTo(ClientDataSourceLink::class.java).one()
}

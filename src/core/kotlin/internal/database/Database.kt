package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.DatabaseConnection

internal fun getDatabaseConnectionInfoFromId(
    hd : Handle,
    dbId : Long
) : DatabaseConnection {
    return hd.select("""
        SELECT *
        FROM database_connection_info
        WHERE db_id = ?
    """, dbId).mapTo(DatabaseConnection::class.java).one()
}

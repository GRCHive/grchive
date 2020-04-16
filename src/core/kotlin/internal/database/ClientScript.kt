package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.ClientScript

internal fun getClientScriptFromCodeLinkId(
    hd : Handle,
    id : Long
) : ClientScript {
    return hd.select("""
        SELECT cs.*
        FROM client_scripts AS cs
        INNER JOIN code_to_client_scripts_link AS link
            ON link.script_id = cs.id
        WHERE link.id = ?
    """, id).mapTo(ClientScript::class.java).one()
}

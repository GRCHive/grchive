package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.SupportedParameterType

internal fun getSupportedParameterTypeFromId(
    hd : Handle,
    id : Int
) : SupportedParameterType {
    return hd.select("""
        SELECT *
        FROM supported_code_parameter_types
        WHERE id = ?
    """, id).mapTo(SupportedParameterType::class.java).one()
}

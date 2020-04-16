package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * A client-defined script object.
 *
 * @property id Unique database ID.
 * @property orgId Organization ID that this belongs to.
 * @property name Human readable name of the data object.
 * @property description Description of this data object.
 */
data class ClientScript (
    @ColumnName("id") val id : Long,
    @ColumnName("org_id") val orgId : Int,
    @ColumnName("name") val name: String,
    @ColumnName("description") val description: String
)

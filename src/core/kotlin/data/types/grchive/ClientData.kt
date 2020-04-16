package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * A client-defined data object.
 *
 * @property id Unique database ID.
 * @property orgId Organization ID that this belongs to.
 * @property name Human readable name of the data object.
 * @property description Description of this data object.
 */
data class ClientData (
    @ColumnName("id") val id : Long,
    @ColumnName("org_id") val orgId : Int,
    @ColumnName("name") val name: String,
    @ColumnName("description") val description: String
)

/**
 * A potential source of data made available to our clients.
 *
 * @property id Unique database ID.
 * @property name Human readable name.
 * @property kotlinClass Full path to the Kotlin class to create.
 */
data class DataSourceOption (
    @ColumnName("id") val id : Long,
    @ColumnName("name") val name : String,
    @ColumnName("kotlin_class") val kotlinClass : String
)

/**
 * These values should match up to what's represented by the id value of [DataSourceOption].
 */
enum class SupportedDataSources(val id : Long) {
    kGrchive(1),
    kPostgres(2)
}

/**
 * Defines what each client data is linked to as its "data source." A data source
 * could be GRCHive, their own database, or a cloud vendor.
 *
 * @property orgId Organization ID that this belongs to.
 * @property dataId ID of the [ClientData] this is linked to.
 * @property sourceId ID of the [DataSourceOption] the [ClientData] is linked to.
 * @property sourceTarget Additional information on what we actually linked to.
 */
 data class ClientDataSourceLink (
    @ColumnName("org_id") val orgId : Int,
    @ColumnName("data_id") val dataId : Long,
    @ColumnName("source_id") val sourceId : Long,
    @ColumnName("source_target") val sourceTarget : Map<String, Any?>
)

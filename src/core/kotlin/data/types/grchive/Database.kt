package grchive.core.data.types.grchive

import kotlin.text.StringBuilder

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * A client-defined database object.
 *
 */
data class Database (
    @ColumnName("id") val id : Long,
    @ColumnName("name") val name : String,
    @ColumnName("org_id") val orgId : Int,
    @ColumnName("type_id") val typeId : Int,
    @ColumnName("other_type") val otherType : String,
    @ColumnName("version") val version : String
)

data class DatabaseConnection (
    @ColumnName("id") val id : Long,
    @ColumnName("db_id") val dbId : Long,
    @ColumnName("org_id") val orgId : Int,
    @ColumnName("host") val host : String,
    @ColumnName("port") val port : Int,
    @ColumnName("dbname") val dbName : String,
    @ColumnName("parameters") val parameters : Map<String, Any?>,
    @ColumnName("username") val username : String,
    @ColumnName("password") val password : String,
    @ColumnName("salt") val salt : String
)

fun createJdbcConnectionString(conn : DatabaseConnection) : String {
    var info = StringBuilder("${conn.host}:${conn.port}/${conn.dbName}?")
    conn.parameters.forEach {
        k, v ->
            info.append("${k}=${v}&")
    }

    return info.toString()
}

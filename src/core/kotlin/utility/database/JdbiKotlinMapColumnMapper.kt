package grchive.core.utility.database

import org.jdbi.v3.core.mapper.ColumnMapper
import org.jdbi.v3.core.statement.StatementContext

import java.sql.ResultSet

class JdbiKotlinMapColumnMapper : ColumnMapper<Map<String, Any?>> {
    override fun map(r : ResultSet, colIdx : Int, ctx : StatementContext) : Map<String, Any?> {
        var result = mutableMapOf<String, Any?>()
        return result
    }
}

package grchive.core.utility.database

import com.fasterxml.jackson.module.kotlin.*

import org.jdbi.v3.core.mapper.ColumnMapper
import org.jdbi.v3.core.statement.StatementContext

import java.sql.ResultSet

class JdbiKotlinMapColumnMapper : ColumnMapper<Map<String, Any?>> {
    override fun map(r : ResultSet, colIdx : Int, ctx : StatementContext) : Map<String, Any?> {
        val rawJson = r.getString(colIdx)
        val mapper = jacksonObjectMapper()
        return mapper.readValue(rawJson)
    }
}

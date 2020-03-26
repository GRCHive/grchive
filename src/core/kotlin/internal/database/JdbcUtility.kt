package grchive.core.internal.database

import com.zaxxer.hikari.*
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.postgres.PostgresPlugin

import grchive.core.data.types.grchive.*
import grchive.core.internal.DatabaseConfig

internal fun createGrchiveHikariDataSource(cfg : DatabaseConfig) : HikariDataSource {
    val ds = HikariDataSource()
    ds.setJdbcUrl("jdbc:postgresql://${cfg.connection}") 
    ds.setUsername(cfg.username)
    ds.setPassword(cfg.password)
    return ds
}

fun setupGrchiveJdbi(jdbi : Jdbi) {
    jdbi.installPlugin(PostgresPlugin())
    jdbi.registerRowMapper(ApiKeyMapper())
}

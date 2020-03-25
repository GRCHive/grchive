package grchive.core.internal.database

import com.zaxxer.hikari.*
import grchive.core.internal.DatabaseConfig

internal fun createGrchiveHikariDataSource(cfg : DatabaseConfig) : HikariDataSource {
    val ds = HikariDataSource()
    ds.setJdbcUrl("jdbc:${cfg.connection}") 
    return ds
}

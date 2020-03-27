package grchive.core.internal.database

import com.zaxxer.hikari.*
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.mapper.reflect.ConstructorMapper
import org.jdbi.v3.postgres.PostgresPlugin

import grchive.core.data.types.grchive.*
import grchive.core.internal.DatabaseConfig

internal fun createGrchiveHikariDataSource(cfg : DatabaseConfig, ro : Boolean) : HikariDataSource {
    val ds = HikariDataSource()
    ds.setJdbcUrl("jdbc:postgresql://${cfg.connection}") 
    ds.setUsername(cfg.username)
    ds.setPassword(cfg.password)
    ds.setReadOnly(ro)
    return ds
}

fun setupGrchiveJdbi(jdbi : Jdbi) {
    jdbi.registerRowMapper(ConstructorMapper.factory(ApiKey::class.java))
    jdbi.registerRowMapper(ConstructorMapper.factory(Role::class.java))
    jdbi.registerRowMapper(ConstructorMapper.factory(RolePermissions::class.java))
    jdbi.registerRowMapper(ConstructorMapper.factory(FullRole::class.java))
    jdbi.registerRowMapper(ConstructorMapper.factory(ProcessFlowMetadata::class.java))
}

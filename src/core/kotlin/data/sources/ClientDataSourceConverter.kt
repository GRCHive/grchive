package grchive.core.data.sources

import org.jdbi.v3.core.Handle
import com.zaxxer.hikari.*

import grchive.core.api.vault.VaultClient

import grchive.core.data.sources.databases.PostgresDataSource
import grchive.core.data.types.grchive.ApiKey
import grchive.core.data.types.grchive.ClientDataSourceLink
import grchive.core.data.types.grchive.createJdbcConnectionString
import grchive.core.data.types.grchive.SupportedDataSources

import grchive.core.internal.Config
import grchive.core.internal.database.getDatabaseConnectionInfoFromId

import grchive.core.security.decryptPassword

import javax.sql.DataSource

internal fun createGrchiveDataSource(cfg : Config, userId : Long, orgId : Int) : GrchiveDataSource {
    return GrchiveDataSource(cfg, userId, orgId)
}

internal fun createPostgresDataSource(hd : Handle, link : ClientDataSourceLink, vault : VaultClient) : PostgresDataSource {
    val dbId : Long = (link.sourceTarget.get("id")!! as Int).toLong()
    val conn = getDatabaseConnectionInfoFromId(hd, dbId)

    val ds = HikariDataSource()
    ds.setJdbcUrl("jdbc:postgresql://${createJdbcConnectionString(conn)}") 
    ds.setUsername(conn.username)
    ds.setPassword(decryptPassword(conn.password, conn.salt, vault))
    ds.setReadOnly(true)

    return PostgresDataSource(ds)
}

internal fun makeDataSourceFromClientDataSourceLink(
    link : ClientDataSourceLink,
    cfg : Config,
    userId : Long,
    orgId : Int,
    hd : Handle,
    vault : VaultClient
) : RawDataSource {
    return when (link.sourceId) {
        SupportedDataSources.kGrchive.id -> createGrchiveDataSource(cfg, userId, orgId)
        SupportedDataSources.kPostgres.id -> createPostgresDataSource(hd, link, vault)
        else -> {
            throw Exception("Unsupported data source.")
        }
    }
}

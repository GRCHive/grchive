package grchive.core.data.types.grchive

import java.sql.ResultSet
import java.time.OffsetDateTime

import org.jdbi.v3.core.mapper.RowMapper
import org.jdbi.v3.core.statement.StatementContext

import grchive.core.security.hashStringSHA512


/**
 * A GRCHive API key which is used to get access to GRCHive resources.
 *
 * The API key is either linked to a user who has a [Role] or is directly
 * linked to a [Role]. Linking directly to a [Role] takes precedence over
 * linking to a [User].
 *
 * @property id The underlying unique identifier.
 * @property hashedKey The SHA512-hashed key.
 * @property expirationDate When this API key will expire.
 */
internal data class ApiKey (
    val id : Long,
    val hashedKey : String,
    val expirationDate : OffsetDateTime
)

internal class ApiKeyMapper : RowMapper<ApiKey> {
    override fun map(rs : ResultSet, ctx : StatementContext) : ApiKey {
        return ApiKey(
            rs.getLong("id"),
            rs.getString("hashed_api_key"),
            rs.getObject("expiration_date", OffsetDateTime::class.java))
    }
}

fun hashRawApiKey(key : String) : String? {
    return hashStringSHA512(key)
}

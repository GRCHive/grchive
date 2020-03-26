package grchive.core.data.types.grchive

import java.sql.ResultSet
import java.time.OffsetDateTime

import org.jdbi.v3.core.mapper.reflect.ColumnName

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
    @ColumnName("id") val id : Long,
    @ColumnName("hashed_api_key") val hashedKey : String,
    @ColumnName("expiration_date") val expirationDate : OffsetDateTime
)

fun hashRawApiKey(key : String) : String? {
    return hashStringSHA512(key)
}

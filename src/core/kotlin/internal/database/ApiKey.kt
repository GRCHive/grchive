package grchive.core.internal.database

import org.jdbi.v3.core.Handle
import java.util.Optional
import grchive.core.data.types.grchive.ApiKey
import grchive.core.data.types.grchive.hashRawApiKey

/**
 * Finds the [ApiKey] with the given raw (un-hashed) API key.
 *
 * @param hd A JDBI handle.
 * @param rawKey The unhashed key.
 * @return An [ApiKey] if found. Null if not.
 */
internal fun getApiKeyFromRawKey(hd : Handle, rawKey : String) : ApiKey? {
    val res : Optional<ApiKey> = hd.select("""
		SELECT key.*
		FROM api_keys AS key
		WHERE hashed_api_key = ?
    """, hashRawApiKey(rawKey)).mapTo(ApiKey::class.java).findOne()
    return if (res.isPresent()) res.get() else null
}

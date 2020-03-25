package grchive.core.internal.database

import org.jdbi.v3.core.Handle
import java.util.Optional
import grchive.core.data.types.grchive.FullRole

/**
 * Finds the [FullRole] for the given unhashed API key.
 *
 * @param hd A JDBI handle.
 * @param apiKeyId The unique ID of the API key.
 * @return A [FullRole] if found. Null if not.
 */
fun getRoleAttachedToApiKey(hd : Handle, apiKeyId : Long) : FullRole? {
    return null
}

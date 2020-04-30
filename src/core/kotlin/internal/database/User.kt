package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.User

/**
 * Finds the [User]s in a given organization.
 *
 * @param hd A JDBI handle.
 * @param orgId Organization ID to find userse in.
 * @return An [ApiKey] if found. Null if not.
 */
internal fun getUsersInOrganizationId(hd : Handle, orgId: Int) : List<User> {
    val res : List<User> = hd.select("""
        SELECT u.*
        FROM users AS u
        INNER JOIN user_orgs AS uo
            ON uo.user_id = u.id
        WHERE uo.org_id = ?
    """, orgId).mapTo(User::class.java).list()
    return res
}

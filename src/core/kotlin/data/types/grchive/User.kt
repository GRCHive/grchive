package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * Contains basic information about a registered user.
 *
 * @property id Unique database ID.
 * @property firstName
 * @property lastName
 * @property email
 */
data class User(
    @ColumnName("id") val id : Long,
    @ColumnName("first_name") val firstName : String,
    @ColumnName("last_name") val lastName : String,
    @ColumnName("email") val email : String
)

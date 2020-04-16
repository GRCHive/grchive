package grchive.core.data.types.grchive

import org.jdbi.v3.core.mapper.reflect.ColumnName

/**
 * Contains information about the identifier and type of a parameter that will
 * be passed to a script.
 *
 * @property name A human readable name.
 * @property paramId The ID of the parameter type stored in the database.
 */
data class ScriptParameter (
    val name : String,
    val paramId : Int
)

/**
 * Information about a type that we support for letting users input arbitrary values for on the web interface
 * and having that value passed to the script.
 *
 * @property id Unique database ID.
 * @property name Human-readable name of the type.
 * @property golangType The type to use in Golang.
 * @property kotlinType The type to use in Kotlin.
 */
data class SupportedParameterType (
    @ColumnName("id") val id : Int,
    @ColumnName("name") val name : String,
    @ColumnName("golang_type") val golangType : String,
    @ColumnName("kotlin_type") val kotlinType : String
)

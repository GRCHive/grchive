package grchive.core.internal

import java.nio.file.*
import org.tomlj.*

/**
 * Configuration of how GRCHive apps connect to the main database.
 *
 * @param connection A valid JDBC connection string without the jdbc: prefix.
 */

internal data class DatabaseConfig (
    val connection : String
)

/**
 * Configuration of GRCHive web apps shared across all applications.
 *
 * @param database Configuration for connecting to the GRCHive database.
 */
class Config {
    internal val database : DatabaseConfig

    constructor(filename : String) {
        val filenamePath : Path = Paths.get(filename)
        if (Files.notExists(filenamePath)) {
            error("Can not find config file at [$filename].")
        }

        val result : TomlParseResult = Toml.parse(filenamePath)
        if (result.hasErrors()) {
            result.errors().forEach {
                System.err.println(it.toString())
            }
            error("Failed to parse config file [$filename].")
        }

        database = DatabaseConfig(
            result.getString("database.connection")!!
        )
    }
}

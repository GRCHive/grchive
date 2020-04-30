package grchive.core.internal

import java.nio.file.*
import org.tomlj.*

import grchive.core.api.vault.VaultConfig

/**
 * Configuration of how GRCHive apps connect to the main database.
 *
 * @param connection A valid JDBC connection string without the jdbc:postgres:// prefix.
 * @param username The username to connect to the DB with.
 * @param password The password to connect to the DB with.
 */

internal data class DatabaseConfig (
    val connection : String,
    val username : String,
    val password : String
)

/**
 * Configuration of GRCHive web apps shared across all applications.
 *
 * @param database Configuration for connecting to the GRCHive database.
 */
class Config {
    internal val database : DatabaseConfig
    internal val vault : VaultConfig

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
            result.getString("database.connection")!!,
            result.getString("database.username")!!,
            result.getString("database.password")!!
        )

        vault = VaultConfig(
            result.getString("vault.url")!!,
            result.getString("vault.userpass.username")!!,
            result.getString("vault.userpass.password")!!
        )
    }

    internal constructor(
        db : DatabaseConfig,
        v : VaultConfig
    ) {
        database = db
        vault = v
    }
}

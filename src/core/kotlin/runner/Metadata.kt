package grchive.core.runner

import java.io.InputStream

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory
import com.fasterxml.jackson.module.kotlin.KotlinModule

import grchive.core.data.types.grchive.ScriptParameter

/**
 * Contains information about what we need to create and pass to the client script.
 *
 * @property params These are the parameters that the user's function is expecting to have.
 * @property clientDataId These are the data sources that the user's function is expecting to be able to use.
 */
data class Metadata(val params : Array<ScriptParameter>, val clientDataId : Array<Long>)

/**
 * Loads the [Metadata] data class from the given resource path.
 *
 * @param resource The resource path in the JAR.
 */
fun loadMetadataFromResource(resource : String) : Metadata {
    val yamlStream : InputStream? = Metadata::class.java.getResourceAsStream(resource)
    if (yamlStream == null) {
        throw Exception("Failed to find script metadata.")
    }

    val mapper = ObjectMapper(YAMLFactory())
    mapper.registerModule(KotlinModule())
    return mapper.readValue(yamlStream, Metadata::class.java)
}

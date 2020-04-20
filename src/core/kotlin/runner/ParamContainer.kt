package grchive.core.runner

import grchive.core.internal.database.getSupportedParameterTypeFromId
import grchive.core.internal.database.getScriptRunParameterValues

import org.jdbi.v3.core.Handle

/**
 * A container that stores key, value pairs of user-defined identifiers and
 * properly typed parameters that need to be passed from our web interface to the scripts.
 */
class ParamContainer {
    private val paramTypes : MutableMap<String, String> = mutableMapOf<String, String>()
    private val paramValues : MutableMap<String, Any?> = mutableMapOf<String, Any?>()

    internal fun addParamType(k : String, typ : String, v : Any?) {
        paramTypes.put(k, typ)
        paramValues.put(k, v)
    }

    fun keys() : Set<String> {
        return paramTypes.keys
    }


    fun getType(k : String) : String? {
        val typ = paramTypes.get(k)
        if (typ == null) {
            return null
        }

        return typ
    }

    fun getValue(k : String) : Any? {
        val v = paramValues.get(k)
        if (v == null) {
            return null
        }

        return v
    }
}

/**
 * Creates a [ParamContainer].
 *
 * @return A [ParamContainer] that knows the types and values of each user-provided parameter.
 * @param handle A JDBI handle to connect to the GRCHive database.
 * @param meta A [Metadata] object that holds what parameters we need to load.
 */
internal fun loadParamContainer(runId : Long, handle : Handle, meta : Metadata) : ParamContainer {
    var container = ParamContainer()
    val paramValues = getScriptRunParameterValues(handle, runId)

    meta.params.forEach {
        val paramType = getSupportedParameterTypeFromId(handle, it.paramId)
        // TODO: This needs to get converted to the proper value.
        container.addParamType(it.name, paramType.name, paramValues.get(it.name))
    }

    return container
}

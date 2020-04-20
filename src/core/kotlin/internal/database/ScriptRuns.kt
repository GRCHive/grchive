package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.ScriptRun
import grchive.core.data.types.grchive.ScriptRunParameter

internal fun getScriptRunFromId(hd : Handle, id: Long) : ScriptRun {
    return hd.select("""
        SELECT run.*
        FROM script_runs AS run
        WHERE run.id = ?
    """, id).mapTo(ScriptRun::class.java).one()
}

internal fun getScriptRunParameterValues(hd : Handle, runId: Long) : Map<String, Any?> {
    val params = hd.select("""
        SELECT *
        FROM script_run_parameters
        WHERE run_id = ?
    """, runId)
        .mapTo(ScriptRunParameter::class.java).list()

    var result = mutableMapOf<String, Any?>()
    params.forEach {
        result.put(it.paramName, it.vals)
    }
    return result
}

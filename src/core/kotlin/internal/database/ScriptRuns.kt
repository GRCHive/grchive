package grchive.core.internal.database

import org.jdbi.v3.core.Handle

import grchive.core.data.types.grchive.ScriptRun

internal fun getScriptRunFromId(hd : Handle, id: Long) : ScriptRun {
    return hd.select("""
        SELECT run.*
        FROM script_runs AS run
        WHERE run.id = ?
    """, id).mapTo(ScriptRun::class.java).one()
}

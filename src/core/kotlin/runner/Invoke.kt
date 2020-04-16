package grchive.core.runner

import grchive.core.internal.Config
import grchive.core.internal.database.*

import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.kotlin.*

/**
 * Runs the specified static function in the class with the given metadata.
 * This function assumes that this function is of standard form written in a way
 * that's expected to be used in our automation process.
 *
 * @param cls The Java class name that wraps the function to run.
 * @param fn The function name to run.
 * @param meta The [Metadata] which we use to get parameters and data sources to pass to the function.
 */
fun invokeWithMetadata(runId : Long, cls : String, fn : String, meta : Metadata) {
    // Load in config from disk. This requires us to know the config path -- hard code this for now...
    val cfg = Config("/config/config.toml")

    // Create a JDBI connection for internal use - this should not get passed to the client for
    // whatever reason since it's not read only.
    val ds = createGrchiveHikariDataSource(cfg.database, false)
    val jdbi = Jdbi.create(ds)
    setupGrchiveJdbi(jdbi)

    // Get the current run status to determine what user is requesting this run.
    // Generate a temporary API key for this user.
    val scriptRun = jdbi.withHandleUnchecked {
        getScriptRunFromId(it, runId)
    }

    // Need to also determine the org ID for which we're running this script.
    val orgId = jdbi.withHandleUnchecked {
        val script = getClientScriptFromCodeLinkId(it, scriptRun.linkId)
        script.orgId
    }

    val params = jdbi.withHandleUnchecked {
        loadParamContainer(it, meta)
    }

    val dataSources = jdbi.withHandleUnchecked {
        loadDataSourceContainer(it, meta, cfg, scriptRun.userId, orgId)
    }

    val jClass = Class.forName(cls)
    jClass.getMethod(fn, ParamContainer::class.java, DataSourceContainer::class.java).invoke(null, params, dataSources)
}

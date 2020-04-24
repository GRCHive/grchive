package grchive.core.data.types.grchive

import java.sql.ResultSet
import java.time.OffsetDateTime
import java.util.Properties

import org.jdbi.v3.core.mapper.ColumnMapper
import org.jdbi.v3.core.mapper.reflect.ColumnName
import org.jdbi.v3.core.statement.StatementContext

/**
 * Information about the progress of a client's run request..
 *
 * @property id Unique database ID.
 * @property linkId A unique identifer to be able to find what client script and managed code was used.
 * @property startTime When the request was sent to run this script.
 * @property buildStartTime When the build started.
 * @property buildFinishTime When the build finished.
 * @property buildSuccess Whether the build succeeded.
 * @property runStartTime When the run started.
 * @property runFinishTime When the run finished.
 * @property runSuccess Whether the run succeeded.
 * @property buildLog Vault encrypted logs from the build stage.
 * @property runLog Vault encrypted logs from the run stage.
 * @property userId User ID of the user who requested this run.
 */
data class ScriptRun (
    @ColumnName("id") val id : Long,
    @ColumnName("link_id") val linkId : Long,
    @ColumnName("start_time") val startTime : OffsetDateTime,
    @ColumnName("build_start_time") val buildStartTime : OffsetDateTime?,
    @ColumnName("build_finish_time") val buildFinishTime : OffsetDateTime?,
    @ColumnName("build_success") val buildSuccess : Boolean,
    @ColumnName("run_start_time") val runStartTime : OffsetDateTime?,
    @ColumnName("run_finish_time") val runFinishTime : OffsetDateTime?,
    @ColumnName("run_success") val runSuccess : Boolean,
    @ColumnName("build_log") val buildLog : String?,
    @ColumnName("run_log") val runLog : String?,
    @ColumnName("user_id") val userId : Long
)

/**
 * Holds the actual value to use for a script parameter for a given script run.
 *
 * @property runId The database ID of the script run.
 * @property paramId  The parameter name to use.
 * @property vals The value to use.
 */ 
data class ScriptRunParameter (
    @ColumnName("run_id") val runId : Long,
    @ColumnName("param_name") val paramName : String,
    @ColumnName("vals") val vals : String
)

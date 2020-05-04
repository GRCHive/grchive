package grchive.core.data.track

import org.jdbi.v3.core.statement.SqlLogger
import org.jdbi.v3.core.statement.StatementContext

class TrackedSource(internal val grchiveDataId : Long) {
    internal val childData = arrayListOf<TrackedData<*>>()
    internal var src : String = ""

    internal fun addData(d : TrackedData<*>) {
        childData.add(d)
    }
}

class TrackedSourceLogger(val src : TrackedSource) : SqlLogger {
    override fun logAfterExecution(context : StatementContext) {
        src.src = """
            |SQL: ${context.getRenderedSql().trim()}
            |Bindings: ${context.getBinding().toString().trim()}
        """.trimMargin()
    }
}

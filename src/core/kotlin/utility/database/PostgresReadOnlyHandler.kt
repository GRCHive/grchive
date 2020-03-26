package grchive.core.utility.database

import org.jdbi.v3.core.Handle

final class PostgresReadOnlyHandler : ReadOnlyHandler {
    override fun setHandleReadOnly(hd : Handle) {
        hd.setReadOnly(true)
        hd.createUpdate("""
            SET SESSION CHARACTERISTICS AS TRANSACTION READ ONLY
        """).execute()
    }
}

package grchive.core.utility.database

import org.jdbi.v3.core.Handle

interface ReadOnlyHandler {
    fun setHandleReadOnly(hd : Handle)
}

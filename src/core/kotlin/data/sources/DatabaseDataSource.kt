package grchive.core.data.sources

import javax.sql.DataSource
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.Handle
import org.jdbi.v3.core.kotlin.*

import grchive.core.utility.database.ReadOnlyHandler

/** 
 * An abstraction over the data stored in any database that can be 
 * connected to using JDBC.
 * 
 * @property ds A Java DataSource from which to create connections.
 * @property ro A read only handler for the database in question.
 */
internal open class DatabaseDataSource(val ds : DataSource, val ro : ReadOnlyHandler) : RawDataSource {
    internal val jdbi : Jdbi = Jdbi.create(ds)

    open fun <T> withHandle(cb : (hd : Handle) -> T) : T {
        return jdbi.withHandleUnchecked {
            ro.setHandleReadOnly(it)
            cb(it)
        }
    }
}

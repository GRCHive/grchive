package grchive.core.data.sources

import javax.sql.DataSource
import org.jdbi.v3.core.Jdbi

/** 
 * An abstraction over the data stored in any database that can be 
 * connected to using JDBC.
 */
open class DatabaseDataSource (val ds : DataSource) : RawDataSource {
    internal val jdbi : Jdbi = Jdbi.create(ds)
}

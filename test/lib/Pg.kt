package test.grchive

import org.testcontainers.containers.*
import org.jdbi.v3.core.Jdbi
import org.jdbi.v3.core.Handle
import io.kotest.core.listeners.TestListener
import io.kotest.core.spec.Spec
import com.zaxxer.hikari.*
import org.jdbi.v3.core.kotlin.*
import org.flywaydb.core.Flyway

import grchive.core.internal.database.setupGrchiveJdbi

class GPostgreSQLContainer : PostgreSQLContainer<GPostgreSQLContainer>("postgres:12.2")

open class KotestGrchivePgContainer(val initFn : (handle: Handle) -> Unit): TestListener {
    public var pg : GPostgreSQLContainer = GPostgreSQLContainer()
    private var ds : HikariDataSource? = null
    private var jdbi : Jdbi? = null

    override suspend fun beforeSpec(spec : Spec) {
        pg.start()

        // Migrate the database to have the proper schema.
        // Note: Don't try to use classpath location here since 
        // it's broken when not using JDK 8.
        val flyway = Flyway.configure().dataSource(
            pg.getJdbcUrl(),
            pg.getUsername(),
            pg.getPassword())
            .locations("filesystem:devops/database/webserver/sql")
            .load()
        flyway.migrate()

        // Allow the test to insert data into the database.
        ds = HikariDataSource()
        ds?.setJdbcUrl(pg.getJdbcUrl())
        ds?.setUsername(pg.getUsername())
        ds?.setPassword(pg.getPassword())

        jdbi = Jdbi.create(ds!!)
        setupGrchiveJdbi(jdbi!!)

        jdbi?.useHandleUnchecked {
            initFn(it)
        }
    }

    override suspend fun afterSpec(spec : Spec) {
        pg.stop()
    }

    fun useHandle(fn : (handle : Handle) -> Unit) {
        jdbi!!.useHandleUnchecked(fn)
    }
}

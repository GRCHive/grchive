package grchive.core.test

import org.jdbi.v3.core.Handle

import grchive.core.data.track.TrackedData
import grchive.core.data.track.TrackedSource

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.KotlinModule

enum class TestAction {
    Equal
}

internal data class TrackedTest(
    val action : TestAction,
    val a : TrackedData<*>?,
    val b : TrackedData<*>?,
    val ok : Boolean,
    val field : String
)

class TestContainer {
    private var allTests : ArrayList<TrackedTest> = arrayListOf<TrackedTest>()
    private var uniqueSources : MutableSet<TrackedSource> = mutableSetOf<TrackedSource>()

    internal fun <T1, T2> logTest(action : TestAction, a : TrackedData<T1>?, b : TrackedData<T2>?, ok : Boolean, field : String = "") {
        allTests.add(TrackedTest(action, a, b, ok, field))

        if (a != null) {
            uniqueSources.add(a.source)
        }

        if (b != null) {
            uniqueSources.add(b.source)
        }
    }

    internal fun commit(hd : Handle, runId : Long, orgId : Int) {
        val trackedDataToDbId : MutableMap<TrackedData<*>, Long> = mutableMapOf<TrackedData<*>, Long>()

        val mapper = ObjectMapper()
        mapper.registerModule(KotlinModule())

        uniqueSources.forEach {
            val srcId : Long = hd.createQuery("""
                INSERT INTO test_sources (run_id, data_id, org_id, src)
                VALUES (?, ?, ?, ?)
                RETURNING id
            """)
                .bind(0, runId)
                .bind(1, it.grchiveDataId)
                .bind(2, orgId)
                .bind(3, it.src)
                .mapTo(Long::class.java)
                .one()

            it.childData.forEach {
                dt ->
                    val dtId : Long = hd.createQuery("""
                        INSERT INTO test_data (source_id, data)
                        VALUES (?, ?::jsonb)
                        RETURNING id
                    """)
                        .bind(0, srcId)
                        .bind(1, mapper.writeValueAsString(dt.t))
                        .mapTo(Long::class.java)
                        .one()
                    trackedDataToDbId.put(dt, dtId)
            }
        }

        allTests.forEach {
            hd.execute("""
                INSERT INTO test_tests (data_a_id, data_b_id, ok, action, field)
                VALUES (?, ?, ?, ?, ?)
            """,
                trackedDataToDbId.get(it.a),
                trackedDataToDbId.get(it.b),
                it.ok,
                it.action.name,
                it.field
            )
        }
    }
}

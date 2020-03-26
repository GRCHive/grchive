package test.grchive.core.internal.database

import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe
import io.kotest.assertions.throwables.shouldThrow

import test.grchive.KotestGrchivePgContainer

import grchive.core.data.filters.Filter
import grchive.core.data.types.grchive.ProcessFlowMetadata
import grchive.core.data.types.grchive.FullRole
import grchive.core.data.types.grchive.Role
import grchive.core.data.types.grchive.RolePermissions
import grchive.core.data.types.grchive.emptyRolePermissions
import grchive.core.data.types.grchive.fullRolePermissions
import grchive.core.data.types.grchive.AccessType
import grchive.core.data.types.grchive.ResourcePermissionDeniedException
import grchive.core.data.types.grchive.Resources

import grchive.core.data.filters.*

import grchive.core.internal.database.getAllProcessFlowMetadata

class GetAllProcessFlowMetadataTest: StringSpec({
    var refOrgId : Int = -1
    var refFlows : List<ProcessFlowMetadata> = emptyList<ProcessFlowMetadata>()

    val pg = KotestGrchivePgContainer {
        refOrgId = it.createQuery("""
            INSERT INTO organizations (org_group_id, org_group_name, org_name)
            VALUES ('Blah', 'Blah', 'Blah')
            RETURNING id
        """).mapTo(Int::class.java).one()

        refFlows = it.createQuery("""
            INSERT INTO process_flows (name, org_id, description, created_time, last_updated_time)
            VALUES 
                ('Test 1', :org_id, 'Test 1', NOW(), NOW()),
                ('Test 2', :org_id, 'Test 2', NOW(), NOW()),
                ('Test 3', :org_id, 'Test 3', NOW(), NOW())
            RETURNING *
        """).bind("org_id", refOrgId).mapTo(ProcessFlowMetadata::class.java).list()
    }
    listener(pg)

    "can not find for orgId" {
        pg.useHandle {
            val metadata = getAllProcessFlowMetadata(it, 2222, FullRole(
                Role(-1, "", "", false, false, -1),
                fullRolePermissions()
            ), emptyMap<String, Filter>())

            metadata.size shouldBe 0
        }
    }

    "bad permissions" {
        pg.useHandle {
            val ex = shouldThrow<ResourcePermissionDeniedException> {
                    getAllProcessFlowMetadata(it, 2222, FullRole(
                    Role(-1, "", "", false, false, -1),
                    fullRolePermissions().copy(
                        processFlowsAccess=(AccessType.All.bit xor AccessType.View.bit)
                    )
                ), emptyMap<String, Filter>())
            }

            ex.res shouldBe Resources.ResourceProcessFlows
            ex.access shouldBe AccessType.View
        }
    }

    "find all" {
        pg.useHandle {
            val metadata = getAllProcessFlowMetadata(it, refOrgId, FullRole(
                Role(-1, "", "", false, false, -1),
                fullRolePermissions()
            ), emptyMap<String, Filter>())

            metadata.size shouldBe 3
            metadata.forEachIndexed { idx, it -> 
                it shouldBe refFlows[idx]
            }
        }
    }

    "id equal filter" {
        pg.useHandle {
            val testIdx = 1
            val metadata = getAllProcessFlowMetadata(it, refOrgId, FullRole(
                Role(-1, "", "", false, false, -1),
                fullRolePermissions()
            ), mapOf(
                "id" to EqualFilter<Long>(refFlows[testIdx].id)
            ))

            metadata.size shouldBe 1
            metadata[0] shouldBe refFlows[testIdx]
        }
    }
})

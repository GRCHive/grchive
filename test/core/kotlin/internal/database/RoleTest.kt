package test.grchive.core.internal.database

import io.kotest.matchers.shouldBe
import io.kotest.matchers.nulls.shouldBeNull
import io.kotest.matchers.nulls.shouldNotBeNull

import io.kotest.core.spec.style.StringSpec
import test.grchive.KotestGrchivePgContainer

import grchive.core.data.types.grchive.AccessType
import grchive.core.data.types.grchive.hashRawApiKey
import grchive.core.data.types.grchive.Resources
import grchive.core.data.types.grchive.Role
import grchive.core.data.types.grchive.RolePermissions
import grchive.core.data.types.grchive.getRolePermissionForResource
import grchive.core.data.types.grchive.unionAccessType

import grchive.core.internal.database.resourceToDatabaseMap
import grchive.core.internal.database.resourceToColumnName

class RoleTest: StringSpec({
    val validRawKey = "ABCDEFGHIJKLMNOP"
    var apiKeyId : Long = -1
    var refOrgId : Int = -1
    var refRole : Role = Role(-1, "", "", false, false, -1)
    var refPermissions = RolePermissions(
        unionAccessType(AccessType.None) /* orgUsersAccess */,
        unionAccessType(AccessType.View) /* orgRolesAccess */,
        unionAccessType(AccessType.Edit) /* processFlowsAccess */,
        unionAccessType(AccessType.Manage) /* controlsAccess */,
        unionAccessType(AccessType.View, AccessType.Edit) /* controlDocumentationAccess */,
        unionAccessType(AccessType.Edit, AccessType.Manage) /* controlDocMetadataAccess */,
        unionAccessType(AccessType.Edit) /* risksAccess */,
        unionAccessType(AccessType.None) /* gLAccess */,
        unionAccessType(AccessType.View, AccessType.Manage) /* systemAccess */,
        unionAccessType(AccessType.Edit) /* dbAccess */,
        unionAccessType(AccessType.Manage) /* dbConnectionAccess */,
        unionAccessType(AccessType.View) /* docRequestAccess */,
        unionAccessType(AccessType.Edit) /* deploymentAccess */,
        unionAccessType(AccessType.Manage) /* serverAccess */,
        unionAccessType(AccessType.Manage) /* vendorAccess */,
        unionAccessType(AccessType.Edit) /* dbSqlAccess */,
        unionAccessType(AccessType.View) /* dbSqlQueryAccess */,
        unionAccessType(AccessType.None) /* dbSqlRequestAccess */,
        unionAccessType(AccessType.View) /* clientDataAccess */,
        unionAccessType(AccessType.View, AccessType.Edit, AccessType.Manage) /* managedCodeAccess */,
        unionAccessType(AccessType.Manage) /* clientScriptAccess */,
        unionAccessType(AccessType.Edit, AccessType.Manage) /* scriptRunAccess */,
        unionAccessType(AccessType.Edit) /* buildLogAccess */,
        unionAccessType(AccessType.View, AccessType.Manage) /* shellScriptAccess */
    )

    val pg = KotestGrchivePgContainer {
        apiKeyId = it.createQuery("""
            INSERT INTO api_keys (hashed_api_key, expiration_date)
            VALUES (?, NOW())
            RETURNING id
        """)
            .bind(0, hashRawApiKey(validRawKey))
            .mapTo(Long::class.java)
            .one()

        refOrgId = it.createQuery("""
            INSERT INTO organizations (org_group_id, org_group_name, org_name)
            VALUES ('Blah', 'Blah', 'Blah')
            RETURNING id
        """).mapTo(Int::class.java).one()

        refRole = it.createQuery("""
            INSERT INTO organization_available_roles (
                name,
                description,
                is_default_role,
                is_admin_role,
                org_id
            )
            VALUES (
                'Blah',
                'Blah',
                true,
                false,
                ?
            )
            RETURNING *
        """).bind(0, refOrgId).mapTo(Role::class.java).one()

        enumValues<Resources>().forEach { v -> 
            it.createUpdate("""
                INSERT INTO ${resourceToDatabaseMap[v]} (role_id, org_id, access_type)
                VALUES (?, ?, ?)
            """)
                .bind(0, refRole.id)
                .bind(1, refOrgId)
                .bind(2, getRolePermissionForResource(refPermissions, v))
                .execute()
        }
    }
    listener(pg)
})

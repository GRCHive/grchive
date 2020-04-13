package test.grchive.core.data.sources

import io.kotest.core.spec.style.StringSpec
import io.kotest.matchers.shouldBe
import io.kotest.assertions.throwables.shouldThrow

import grchive.core.data.sources.GrchiveDataSource
import grchive.core.data.types.grchive.AccessType
import grchive.core.data.types.grchive.fullRolePermissions
import grchive.core.data.types.grchive.getRolePermissionForResource
import grchive.core.data.types.grchive.hashRawApiKey
import grchive.core.data.types.grchive.Resources
import grchive.core.data.types.grchive.Role

import grchive.core.internal.Config
import grchive.core.internal.DatabaseConfig
import grchive.core.internal.database.resourceToDatabaseMap

import test.grchive.KotestGrchivePgContainer

class GrchiveDataSourceTest: StringSpec({
    var ds : GrchiveDataSource? = null
    val validRawKey = "ABCDEFGHIJKLMNOP"
    var refRole : Role? = null
    var refOrgId : Int = -1
    var refPv : String = ""

    val pg = KotestGrchivePgContainer {
        val apiKeyId = it.createQuery("""
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
                .bind(0, refRole!!.id)
                .bind(1, refOrgId)
                .bind(2, AccessType.All.bit)
                .execute()
        }

        it.createUpdate("""
            INSERT INTO api_key_roles (role_id, api_key_id, org_id)
            VALUES (?, ?, ?)
        """)
            .bind(0, refRole!!.id)
            .bind(1, apiKeyId)
            .bind(2, refOrgId)
            .execute()

        refPv = it.createQuery("SELECT version()").mapTo(String::class.java).one()
    }
    listener(pg)

    beforeSpec {
        ds = GrchiveDataSource(
            Config(
                DatabaseConfig(
                    pg.ds!!.getJdbcUrl().replace("jdbc:postgresql://", "") + "&readOnly=true",
                    pg.ds!!.getUsername(),
                    pg.ds!!.getPassword()
                )
            ),
            validRawKey,
            refOrgId
        )
    }

    "Check Active Role" {
        ds!!.activeRole.permissions shouldBe fullRolePermissions()
        ds!!.activeRole.role shouldBe refRole!!
    }

    "SELECT OK" {
        ds!!.db.withHandle {
            val testPv : String = it.createQuery("SELECT version()").mapTo(String::class.java).one()
            testPv shouldBe refPv
        }
    }

    "Handle Read Only" {
        ds!!.db.withHandle {
            it.isReadOnly() shouldBe true
            it.getConnection().isReadOnly() shouldBe true
        }
    }

    "INSERT FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    INSERT INTO get_started_interest (name, email)
                    VALUES ('Test', 'test')
                """).execute()
            }
        }
    }

    "DELETE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    DELETE FROM api_keys
                """).execute()
            }
        }
    }

    "UPDATE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    UPDATE organizations
                    SET name = 'Ooogity'
                """).execute()
            }
        }
    }

    "DROP TABLE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    DROP TABLE organizations
                """).execute()
            }
        }
    }

    "CREATE TABLE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    CREATE TABLE test (
                        id BIGSERIAL PRIMARY KEY
                    )
                """).execute()
            }
        }

    }

    "TRUNCATE TABLE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    TRUNCATE TABLE organizations
                """).execute()
            }
        }
    }

    "ALTER TABLE FAIL" {
        shouldThrow<Exception> {
            ds!!.db.withHandle {
                it.createUpdate("""
                    ALTER TABLE organizations
                    ADD COLUMN test BIGINT UNIQUE;
                """).execute()
            }
        }
    }

    "CREATE ROLE FAIL" {

    }

    "DROP ROLE FAIL" {

    }
})
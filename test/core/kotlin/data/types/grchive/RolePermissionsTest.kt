package test.grchive.core.data.types.grchive

import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec

import grchive.core.data.types.grchive.AccessType
import grchive.core.data.types.grchive.unionAccessType
import grchive.core.data.types.grchive.GrchiveResource
import grchive.core.data.types.grchive.Resources
import grchive.core.data.types.grchive.RolePermissions
import grchive.core.data.types.grchive.getRolePermissionForResource

import kotlin.reflect.full.findParameterByName
import kotlin.reflect.full.instanceParameter

import io.kotest.data.forAll
import io.kotest.data.row

class UnionAccessTypeTest: StringSpec({
    "Test" {
        forAll(
            row(arrayOf(AccessType.None), 0),
            row(arrayOf(AccessType.None, AccessType.None), 0),
            row(arrayOf(AccessType.View), 1),
            row(arrayOf(AccessType.View, AccessType.View), 1),
            row(arrayOf(AccessType.Edit), 2),
            row(arrayOf(AccessType.Edit, AccessType.Edit), 2),
            row(arrayOf(AccessType.Manage), 4),
            row(arrayOf(AccessType.Manage, AccessType.Manage), 4),
            row(arrayOf(AccessType.View, AccessType.Edit), 3),
            row(arrayOf(AccessType.View, AccessType.Manage), 5),
            row(arrayOf(AccessType.Edit, AccessType.Manage), 6),
            row(arrayOf(AccessType.View, AccessType.Edit, AccessType.Manage), 7)
        ) {
            inp, ref -> 
                unionAccessType(*inp) shouldBe ref
        }
    }
})

class GetRolePermissionForResourceTest: StringSpec({
    "Test" {
        for (r in enumValues<Resources>()) {
            val p = RolePermissions()
            getRolePermissionForResource(p, r) shouldBe 0

            // This test will 1) find the field that is annotated with the right resource using
            // GrchiveResource and then 2) create a new RolePermissions object that has that
            // field set to a desired test value (non-zero) and 3) check that the value matches
            // up with what getRolePermissionForResource returns.
            for (f in RolePermissions::class.java.getDeclaredFields()) {
                val ann = f.getAnnotation(GrchiveResource::class.java)
                if (ann.res == r) {
                    for (a in 0..7) {
                        val newP = RolePermissions::copy.callBy(
                            mapOf(
                                RolePermissions::copy.instanceParameter!! to p,
                                RolePermissions::copy.findParameterByName(f.getName())!! to a
                            )
                        )
                        getRolePermissionForResource(newP, r) shouldBe a
                    }
                    break
                }
            }
        }

    }
})

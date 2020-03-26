package test.grchive.core.internal.database

import io.kotest.matchers.nulls.shouldBeNull
import io.kotest.matchers.nulls.shouldNotBeNull
import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec

import java.time.OffsetDateTime

import test.grchive.KotestGrchivePgContainer

import grchive.core.internal.database.getApiKeyFromRawKey
import grchive.core.data.types.grchive.hashRawApiKey
import grchive.core.data.types.grchive.ApiKey

class ApiKeyTest: StringSpec({
    val validRawKey = "ABCDEFGHIJKLMNOP"
    var refKey = ApiKey(
        -1,
        hashRawApiKey(validRawKey)!!,
        OffsetDateTime.now()
    )

    val pg = KotestGrchivePgContainer {
         val id = it.createQuery("""
            INSERT INTO api_keys (hashed_api_key, expiration_date)
            VALUES (?, ?)
            RETURNING id
        """)
            .bind(0, refKey.hashedKey)
            .bind(1, refKey.expirationDate)
            .mapTo(Long::class.java)
            .one()
        refKey = ApiKey(id, refKey.hashedKey, refKey.expirationDate)
    }
    listener(pg)

    "getApiKeyFromRawKey - Find Key" {
        pg.useHandle {
            val key = getApiKeyFromRawKey(it, validRawKey)
            key.shouldNotBeNull()

            key shouldBe refKey
        }

    }

    "getApiKeyFromRawKey - No Key" {
        pg.useHandle {
            val key = getApiKeyFromRawKey(it, "Blah")
            key.shouldBeNull()
        }
    }
})

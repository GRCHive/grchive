package test.grchive.core.internal.database

import io.kotest.matchers.nulls.shouldBeNull
import io.kotest.matchers.nulls.shouldNotBeNull

import io.kotest.core.spec.style.StringSpec
import test.grchive.KotestGrchivePgContainer

import grchive.core.internal.database.getApiKeyFromRawKey
import grchive.core.data.types.grchive.hashRawApiKey

class ApiKeyTest: StringSpec({
    val validRawKey = "ABCDEFGHIJKLMNOP"

    val pg = KotestGrchivePgContainer {
        it.createUpdate("""
            INSERT INTO api_keys (hashed_api_key, expiration_date)
            VALUES (?, NOW())
        """).bind(0, hashRawApiKey(validRawKey)).execute()
    }
    listener(pg)

    "getApiKeyFromRawKey - Find Key" {
        pg.useHandle {
            val key = getApiKeyFromRawKey(it, validRawKey)
            key.shouldNotBeNull()
        }

    }

    "getApiKeyFromRawKey - No Key" {
        pg.useHandle {
            val key = getApiKeyFromRawKey(it, "Blah")
            key.shouldBeNull()
        }
    }
})

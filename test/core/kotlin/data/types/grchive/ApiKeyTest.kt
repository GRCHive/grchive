package test.grchive.core.data.types.grchive

import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec

import grchive.core.security.hashStringSHA512
import grchive.core.data.types.grchive.hashRawApiKey

class HashRawApiKeyTest: StringSpec({
    "should be equivalent to hashStringSHA512" {
        for (key in arrayOf(
            "",
            "ABC",
            "123"
        )) {
            hashRawApiKey(key) shouldBe hashStringSHA512(key)
        }
    }
})

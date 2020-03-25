package test.grchive.core.security

import grchive.core.security.hashStringSHA512
import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec
import io.kotest.data.forAll
import io.kotest.data.row

class HexEncodeTest: StringSpec({
    "SHA512" {
        forAll(
            row(
                "aaabbbcccdddeeefffggghhhiiijjjkkklll",
                "c7ba6e0323bd282fc69b0af27f0ec2021613c82042d36745fc5027e02b26f7c57f6de94753db1c89508c8d35a07f8ba1477182adea1c612a633dc9fa727c5b47"
            ),
            row(
                "abcdef",
                "e32ef19623e8ed9d267f657a81944b3d07adbb768518068e88435745564e8d4150a0a703be2a7d88b61e3d390c2bb97e2d4c311fdc69d6b1267f05f59aa920e7"
            )
        ) {
            inp, ref -> 
                hashStringSHA512(inp) shouldBe ref
        }
    }
})

package test.grchive.core.utility

import grchive.core.utility.hexEncode
import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec
import io.kotest.data.forAll
import io.kotest.data.row

class HexEncodeTest: StringSpec({
    "Hex Basic" {
        forAll(
            row(0x0, "00"),
            row(0x1, "01"),
            row(0x2, "02"),
            row(0x3, "03"),
            row(0x4, "04"),
            row(0x5, "05"),
            row(0x6, "06"),
            row(0x7, "07"),
            row(0x8, "08"),
            row(0x9, "09"),
            row(0xA, "0a"),
            row(0xB, "0b"),
            row(0xC, "0c"),
            row(0xD, "0d"),
            row(0xE, "0e"),
            row(0xF, "0f")
        ) {
            inp, ref -> 
                hexEncode(ByteArray(1){ inp.toByte() }) shouldBe ref
        }
    }

    "Hex Longer" {
        forAll(
            row(0x0000, "00000000"),
            row(0x1230FA, "001230fa"),
            row(0xABCDEF, "00abcdef")
        ) {
            inp, ref -> 
                val testByte = ByteArray(4)
                testByte.set(0, 0x00)
                testByte.set(1, ((inp shr 16) and 0xFF).toByte())
                testByte.set(2, ((inp shr 8) and 0xFF).toByte())
                testByte.set(3, (inp and 0xFF).toByte())

                hexEncode(testByte) shouldBe ref
        }

    }
})

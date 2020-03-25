package grchive.core.utility

import java.lang.StringBuilder

const val kHexChars = "0123456789ABCDEF"

/**
 * Hex encodes the input byte array.
 *
 * @param input The byte array to encode.
 * @return A string that is the hex-encoded equivalent of the input byte array.
 */
fun hexEncode(input : ByteArray) : String {
    var builder = StringBuilder(input.size * 2)
    input.forEach {
        builder.append(kHexChars[(it.toInt() shr 4) and 0xF])
        builder.append(kHexChars[it.toInt() and 0xF])
    }
    return builder.toString()
}

package grchive.core.utility

import java.lang.StringBuilder

const val kHexChars = "0123456789abcdef"

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

/**
 * Hex decodes the input string to a byte array.
 *
 * @param input The hex string to decode.
 * @return The byte array represented by the hex string.
 */
fun hexDecode(input : String) : ByteArray {
    // Need to construct this ourselves instead of using the built-in
    // Integer.parseInt so that the decoded byte array ONLY contains
    // the data in the input string.
    var lowerInput = input.toLowerCase()
    var arr = ByteArray(input.length / 2)

    for (idx in 0 until arr.size) {
        arr[idx] = Integer.parseInt(lowerInput.substring(idx * 2, idx * 2 +2), 16).toByte()
    }
    return arr
}


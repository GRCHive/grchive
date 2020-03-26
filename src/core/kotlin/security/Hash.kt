package grchive.core.security

import java.security.MessageDigest
import java.security.NoSuchAlgorithmException

import grchive.core.utility.hexEncode

/**
 * A generic function that is a thin wrapper over Java's MessageDigest to compute the hash of a string.
 *
 * @param alg The hash algorithm to use. This must be an algorithm that java.security.MessageDigest supports.
 * @param inp The input string to hash.
 * @return The hex-encoded hash. Returns null if the algorithm is not supported.
 */
fun genericStringHash(alg : String, inp : String) : String? {
    try {
        val md : MessageDigest = MessageDigest.getInstance(alg)
        md.update(inp.toByteArray())

        val hashData : ByteArray = md.digest()
        return hexEncode(hashData)
    } catch (e : NoSuchAlgorithmException) {
        return null
    }
}

/**
 * Convenience function to compute the hex encoded SHA-512 of a string.
 *
 * @param input The string to hash.
 * @return The hex-encoded SHA-512 hash.
 */
fun hashStringSHA512(inp : String) : String? {
    return genericStringHash("SHA-512", inp)
}

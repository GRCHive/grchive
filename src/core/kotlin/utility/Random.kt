package grchive.core.utility

import java.lang.StringBuilder
import kotlin.random.Random

const val kRandomChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()-_=+/?;:"

/**
 * Generates a random string of the specified length.
 *
 * @param length Length in characters for the output random string.
 * @return A random string.
 */
fun randomString(length : Int) : String {
    var builder = StringBuilder(length)
    for (i in 0 until length) {
        builder.append(kRandomChars[Random.nextInt(0, kRandomChars.length)])
    }
    return builder.toString()
}

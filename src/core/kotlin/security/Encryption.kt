package grchive.core.security

import grchive.core.api.vault.VaultClient
import grchive.core.api.vault.TransitDecryptBody

import grchive.core.utility.hexDecode

internal fun decryptPassword(pw : String, salt : String, vault : VaultClient) : String {
    val decryptedData = vault.transit.decodedDecrypt("passwords", pw)
    val decryptedString = decryptedData.data.plaintext
    val splitData = decryptedString.split(".")

    if (splitData.size == 0) {
        throw Exception("Invalid encrypted password format.")
    }

    val expectedSalt = splitData.getOrNull(1)!!
    if (expectedSalt != salt) {
        throw Exception("Invalid salt.")
    }

    return hexDecode(splitData[0]).toString(Charsets.UTF_8)
}

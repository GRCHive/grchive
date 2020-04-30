package grchive.core.api.vault

import com.fasterxml.jackson.annotation.JsonIgnoreProperties

import java.util.Base64

import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.Path
import retrofit2.http.POST

@JsonIgnoreProperties(ignoreUnknown = true)
data class TransitDecryptDataResponse(
    val plaintext : String
)

@JsonIgnoreProperties(ignoreUnknown = true)
data class TransitDecryptResponse(
    val data : TransitDecryptDataResponse
)

data class TransitDecryptBody (
    val ciphertext : String
)

internal interface HTTPTransitInterface {
    @POST("v1/transit/decrypt/{name}")
    fun decrypt(@Path("name") name : String, @Body body : TransitDecryptBody) : Call<TransitDecryptResponse>
}

class TransitInterface {
    private val http : HTTPTransitInterface

    internal constructor(inHttp : HTTPTransitInterface) {
        http = inHttp
    }

    fun decodedDecrypt(name : String, cipher : String) : TransitDecryptResponse {
        val resp = http.decrypt(name, TransitDecryptBody(cipher)).execute()
        if (!resp.isSuccessful()) {
            throw Exception("Failed to decrypt: ${resp.code()} -- ${resp.errorBody()!!.string()}")
        }

        val decryptedData = resp.body()!!
        return TransitDecryptResponse(
            TransitDecryptDataResponse(
                Base64.getDecoder().decode(decryptedData.data.plaintext).toString(Charsets.UTF_8)
            )
        )
    }
}


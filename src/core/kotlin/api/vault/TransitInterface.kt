package grchive.core.api.vault

import com.fasterxml.jackson.annotation.JsonIgnoreProperties
import com.fasterxml.jackson.annotation.JsonProperty

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

data class TransitEncryptBatchOutput (
    val ciphertext : String
)

@JsonIgnoreProperties(ignoreUnknown = true)
data class TransitEncryptDataResponse(
    val ciphertext : String?,
    @JsonProperty("batch_results") val batchResults : List<TransitEncryptBatchOutput>? 
)

@JsonIgnoreProperties(ignoreUnknown = true)
data class TransitEncryptResponse(
    val data : TransitEncryptDataResponse
)

data class TransitEncryptBatchInput (
    val plaintext : String
)

data class TransitEncryptBody (
    val plaintext : String,
    @JsonProperty("batch_input") val batchInput : List<TransitEncryptBatchInput>? = null
)

internal interface HTTPTransitInterface {
    @POST("v1/transit/decrypt/{name}")
    fun decrypt(@Path("name") name : String, @Body body : TransitDecryptBody) : Call<TransitDecryptResponse>

    @POST("v1/transit/encrypt/{name}")
    fun encrypt(@Path("name") name : String, @Body body : TransitEncryptBody) : Call<TransitEncryptResponse>
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

    fun batchEncrypt(name : String, raw : List<String>) : TransitEncryptResponse {
        val batchInputs = raw.map {
            TransitEncryptBatchInput(
                Base64.getEncoder().encode(it.toByteArray(Charsets.UTF_8)).toString(Charsets.UTF_8)
            )
        }

        val resp = http.encrypt(name, TransitEncryptBody("", batchInputs)).execute()
        if (!resp.isSuccessful()) {
            throw Exception("Failed to encrypt: ${resp.code()} -- ${resp.errorBody()!!.string()}")
        }

        return resp.body()!!
    }
}


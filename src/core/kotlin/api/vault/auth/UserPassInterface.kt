package grchive.core.api.vault.auth

import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.Path
import retrofit2.http.POST

import com.fasterxml.jackson.annotation.JsonProperty
import com.fasterxml.jackson.annotation.JsonIgnoreProperties

@JsonIgnoreProperties(ignoreUnknown = true)
data class AuthResponse (
    @JsonProperty("client_token") val clientToken : String
)

@JsonIgnoreProperties(ignoreUnknown = true)
data class UserPassLoginResponse(
    val auth : AuthResponse
)

data class UserPassLoginBody (
    val password : String
)

interface UserPassInterface {
    @POST("v1/auth/userpass/login/{username}")
    fun login(@Path("username") username : String, @Body body : UserPassLoginBody) : Call<UserPassLoginResponse>
}

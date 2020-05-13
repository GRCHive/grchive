package grchive.core.api.vault

import com.fasterxml.jackson.module.kotlin.*

import grchive.core.api.vault.auth.UserPassInterface
import grchive.core.api.vault.auth.UserPassLoginBody
import grchive.core.utility.tls.okhttpClientTlsInsecure

import okhttp3.OkHttpClient
import okhttp3.Interceptor
import okhttp3.Response

import okio.Buffer

import retrofit2.converter.jackson.JacksonConverterFactory
import retrofit2.Retrofit

data class VaultConfig (
    val url : String,
    val username : String,
    val password : String
)

private class JSONHeaderInterceptor : Interceptor {
    override fun intercept(chain : Interceptor.Chain) : Response {
        val newRequest = chain.request().newBuilder()
            .header("Content-Type", "application/json")
            .build()
        return chain.proceed(newRequest)
    }
}

private class VaultAuthInterceptor(val token : String) : Interceptor {
    override fun intercept(chain : Interceptor.Chain) : Response {
        val newRequest = chain.request().newBuilder()
            .header("X-Vault-Token", token)
            .build()
        return chain.proceed(newRequest)
    }
}

class VaultClient (
    val cfg : VaultConfig
) {
    private val token : String
    private val retrofit : Retrofit

    // HTTP API Interfaces that we expose to the user.
    private val httpTransit : HTTPTransitInterface
    val transit : TransitInterface

    init {
        // Create a separate auth retrofit instance since the main retrofit
        // instance will need to set the authorization header.
        var authClient : OkHttpClient
        
        if (System.getenv("ALLOW_TLS_INSECURE") == null) {
            authClient = OkHttpClient.Builder()
                .addInterceptor(JSONHeaderInterceptor())
                .build()
        } else {
            authClient = okhttpClientTlsInsecure().newBuilder()
                .addInterceptor(JSONHeaderInterceptor())
                .build()
        }

        val authRetrofit = Retrofit.Builder()
            .baseUrl(cfg.url)
            .addConverterFactory(JacksonConverterFactory.create(
                jacksonObjectMapper()
            ))
            .client(authClient)
            .build()

        val userpass = authRetrofit.create(UserPassInterface::class.java)
        val resp = userpass.login(cfg.username, UserPassLoginBody(cfg.password)).execute()
        if (!resp.isSuccessful()) {
            throw Exception("Failed to auth with Vault: ${resp.code()}")
        }
        token = resp.body()!!.auth.clientToken

        // Now add the authorization token to the header using the OkHttpClient and use
        // this Retrofit instance going forward.
        retrofit = authRetrofit.newBuilder().client(
            authClient.newBuilder().addInterceptor(VaultAuthInterceptor(token)).build()
        ).build()
        httpTransit = retrofit.create(HTTPTransitInterface::class.java)
        transit = TransitInterface(httpTransit)
    }
}

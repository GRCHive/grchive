package grchive.core.utility.tls 

import javax.net.ssl.HostnameVerifier
import javax.net.ssl.SSLContext
import javax.net.ssl.SSLSession
import javax.net.ssl.X509TrustManager

import java.security.cert.X509Certificate
import java.security.SecureRandom

import okhttp3.OkHttpClient

internal class InsecureX509TrustManager : X509TrustManager {
    override fun checkClientTrusted(chain : Array<X509Certificate>, authType : String) {
    }

    override fun checkServerTrusted(chain : Array<X509Certificate>, authType : String) {
    }

    override fun getAcceptedIssuers() : Array<X509Certificate> {
        return emptyArray<X509Certificate>()
    }
}

internal class InsecureHostnameVerifier : HostnameVerifier {
    override fun verify(hostname : String, session : SSLSession) : Boolean {
        return true
    }
}

fun okhttpClientTlsInsecure() : OkHttpClient {
    val trustManagers = arrayOf(InsecureX509TrustManager())
    val sslContext = SSLContext.getInstance("TLSv1.2")
    sslContext.init(null, trustManagers, SecureRandom())

    return OkHttpClient.Builder()
        .sslSocketFactory(
            sslContext.getSocketFactory(),
            trustManagers[0]
        )
        .hostnameVerifier(InsecureHostnameVerifier())
        .build()
}

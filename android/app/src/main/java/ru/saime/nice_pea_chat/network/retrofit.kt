package ru.saime.nice_pea_chat.network

import com.skydoves.retrofit.adapters.result.ResultCallAdapterFactory
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.Response
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import java.net.SocketTimeoutException
import java.net.URL


private const val host = "http://192.168.31.94:7511"

fun retrofit(
    npcClientStore: NpcClientStore,
): Retrofit {
    val logging = HttpLoggingInterceptor()
        .setLevel(HttpLoggingInterceptor.Level.BODY)

    val client = OkHttpClient.Builder()
        .addInterceptor(ReplaceNpcUrlPlaceholderInterceptor(npcClientStore))
        .addInterceptor(logging)
        .addInterceptor(RetryInterceptor(3))
        .build()

    return Retrofit.Builder()
        .baseUrl(NpcUrlPlaceholder)
        .client(client)
        .addConverterFactory(GsonConverterFactory.create())
        .addCallAdapterFactory(ResultCallAdapterFactory.create())
        .build()
}


private class RetryInterceptor(private val retryAttempts: Int) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        repeat(retryAttempts) {
            try {
                return chain.proceed(chain.request())
            } catch (e: SocketTimeoutException) {
                e.printStackTrace()
            }
        }
        throw RuntimeException("failed to compile the request")
    }
}


const val NpcProtocolPlaceholder = "http"
const val NpcHostPlaceholder = "<npc_host>"
const val NpcPortPlaceholder = 7511
const val NpcUrlPlaceholder = "$NpcProtocolPlaceholder://$NpcHostPlaceholder:$NpcPortPlaceholder"

fun npcBaseUrl(store: NpcClientStore, default: String = ""): String {
    return if (store.baseUrl != "") {
        store.baseUrl
    } else if (default != "") {
        default
    } else {
        NpcUrlPlaceholder
    }
}

private class ReplaceNpcUrlPlaceholderInterceptor(
    private val store: NpcClientStore,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val urlString = chain.request().url.toString()
        if (urlString.startsWith(NpcUrlPlaceholder) && store.baseUrl != "") {
            val newUrl = URL(store.baseUrl + urlString.removePrefix(NpcUrlPlaceholder))
            val request = chain.request().newBuilder()
                .url(newUrl)
                .build()
            return chain.proceed(request)
        }

        return chain.proceed(chain.request())
    }
}

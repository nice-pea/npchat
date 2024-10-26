package ru.saime.nice_pea_chat.network

import com.skydoves.retrofit.adapters.result.ResultCallAdapterFactory
import okhttp3.HttpUrl
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.Response
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import java.net.SocketTimeoutException


private const val host = "http://192.168.31.94:7511"

fun retrofit(
    npcClientStore: NpcClientStore,
): Retrofit {
    val logging = HttpLoggingInterceptor()
        .setLevel(HttpLoggingInterceptor.Level.BODY)

    val client = OkHttpClient.Builder()
        .addInterceptor(InsertUrlInterceptor(npcClientStore))
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

fun npcUrl(store: NpcClientStore, default: String = ""): String {
    return if (store.urlIsNotEmpty()) {
        "$NpcProtocolPlaceholder://${store.host}:${store.port}"
    } else if (default.isNotBlank()) {
        default
    } else {
        NpcUrlPlaceholder
    }
}

private class InsertUrlInterceptor(
    private val store: NpcClientStore,
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val url = chain.request().url
        if (
            url.host == NpcHostPlaceholder && url.port == NpcPortPlaceholder
            && store.port != 0 && store.host != ""
        ) {
            val newUrl = url.newBuilder()
                .port(store.port)
                .host(store.host)
                .build()
            HttpUrl.Builder()
            val request = chain.request().newBuilder()
                .url(newUrl)
                .build()
            return chain.proceed(request)
        } else {
            return chain.proceed(chain.request())
        }
    }
}

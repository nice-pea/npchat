package ru.saime.nice_pea_chat.di

import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import org.koin.core.module.dsl.singleOf
import org.koin.core.module.dsl.viewModelOf
import org.koin.dsl.module
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import ru.saime.nice_pea_chat.data.api.ApiClientApi
import ru.saime.nice_pea_chat.data.api.AuthenticationApi
import ru.saime.nice_pea_chat.data.repositories.ApiClient
import ru.saime.nice_pea_chat.data.repositories.AuthenticationRepository
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationViewModel
import ru.saime.nice_pea_chat.screens.login.LoginViewModel

private const val host = "http://192.168.31.94:7511"

val appModule = module {
    single {
        val logging = HttpLoggingInterceptor()
            .setLevel(HttpLoggingInterceptor.Level.BODY)

        val client = OkHttpClient.Builder()
            .addInterceptor(logging)
            .build()

        Retrofit.Builder()
            .baseUrl("http://example.com")
            .client(client)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
    }

    single {
        get<Retrofit>().create(AuthenticationApi::class.java)
    }
    single {
        get<Retrofit>().create(ApiClientApi::class.java)
    }

    singleOf(::AuthenticationRepository)
    singleOf(::AuthenticationStore)
    singleOf(::ApiClient)

    viewModelOf(::AuthenticationViewModel)
    viewModelOf(::LoginViewModel)
}
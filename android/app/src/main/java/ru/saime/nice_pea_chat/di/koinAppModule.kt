package ru.saime.nice_pea_chat.di

import org.koin.core.module.dsl.singleOf
import org.koin.core.module.dsl.viewModelOf
import org.koin.dsl.module
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import ru.saime.nice_pea_chat.AuthenticationApi
import ru.saime.nice_pea_chat.AuthenticationRepository
import ru.saime.nice_pea_chat.AuthnStore
import ru.saime.nice_pea_chat.MainViewModel

val appModule = module {
    single {
        Retrofit.Builder()
            .baseUrl("http://192.168.31.94:7511")
            .addConverterFactory(GsonConverterFactory.create())
            .build()
    }
    single {
        get<Retrofit>().create(AuthenticationApi::class.java)
    }
    singleOf(::AuthenticationRepository)
    singleOf(::AuthnStore)

    viewModelOf(::MainViewModel)
}
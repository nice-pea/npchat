package ru.saime.nice_pea_chat.di

import org.koin.core.module.dsl.singleOf
import org.koin.core.module.dsl.viewModelOf
import org.koin.dsl.module
import retrofit2.Retrofit
import ru.saime.nice_pea_chat.data.api.AuthenticationApi
import ru.saime.nice_pea_chat.data.api.NpcClientApi
import ru.saime.nice_pea_chat.data.repositories.AuthenticationRepository
import ru.saime.nice_pea_chat.data.repositories.NpcClient
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import ru.saime.nice_pea_chat.network.retrofit
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationViewModel
import ru.saime.nice_pea_chat.screens.login.LoginViewModel


val appModule = module {
    // LocalStore
    singleOf(::AuthenticationStore)
    singleOf(::NpcClientStore)

    // Api
    single { retrofit(get()) }
    single { get<Retrofit>().create(AuthenticationApi::class.java) }
    single { get<Retrofit>().create(NpcClientApi::class.java) }

    // Repositories
    singleOf(::AuthenticationRepository)
    singleOf(::NpcClient)

    // ViewModels
    viewModelOf(::AuthenticationViewModel)
    viewModelOf(::LoginViewModel)
}
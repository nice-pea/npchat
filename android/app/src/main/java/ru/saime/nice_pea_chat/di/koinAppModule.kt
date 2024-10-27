package ru.saime.nice_pea_chat.di

import org.koin.core.module.dsl.singleOf
import org.koin.core.module.dsl.viewModelOf
import org.koin.core.scope.Scope
import org.koin.dsl.module
import retrofit2.Retrofit
import ru.saime.nice_pea_chat.data.api.AuthenticationApi
import ru.saime.nice_pea_chat.data.api.ChatsApi
import ru.saime.nice_pea_chat.data.api.NpcClientApi
import ru.saime.nice_pea_chat.data.repositories.AuthenticationRepository
import ru.saime.nice_pea_chat.data.repositories.ChatsRepository
import ru.saime.nice_pea_chat.data.repositories.NpcClient
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import ru.saime.nice_pea_chat.network.retrofit
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationViewModel
import ru.saime.nice_pea_chat.screens.chats.ChatsViewModel
import ru.saime.nice_pea_chat.screens.login.LoginViewModel


val appModule = module {
    // LocalStore
    singleOf(::AuthenticationStore)
    singleOf(::NpcClientStore)

    // Api
    single { retrofit(get(), get()) }
    fun <T> Scope.retroApi(service: Class<T>) = get<Retrofit>().create(service)
    single { retroApi(AuthenticationApi::class.java) }
    single { retroApi(NpcClientApi::class.java) }
    single { retroApi(ChatsApi::class.java) }

    // Repositories
    singleOf(::AuthenticationRepository)
    singleOf(::NpcClient)
    singleOf(::ChatsRepository)

    // ViewModels
    viewModelOf(::AuthenticationViewModel)
    viewModelOf(::LoginViewModel)
    viewModelOf(::ChatsViewModel)
}
package ru.saime.nice_pea_chat.di

import org.koin.core.module.dsl.viewModelOf
import org.koin.dsl.module
import ru.saime.nice_pea_chat.MainViewModel

val appModule = module {
    viewModelOf(::MainViewModel)
}
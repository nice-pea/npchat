package ru.saime.nice_pea_chat.common


sealed interface AsyncData<out T> {
    object None : AsyncData<Nothing>
    object Loading : AsyncData<Nothing>
    data class Err(val err: Throwable) : AsyncData<Nothing>
    data class Ok<out T>(val data: T) : AsyncData<T>
}
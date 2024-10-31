package ru.saime.nice_pea_chat.common

import kotlinx.coroutines.flow.MutableStateFlow
import ru.saime.nice_pea_chat.common.AsyncData.None


sealed interface AsyncData<out T> {
    object None : AsyncData<Nothing>
    object Loading : AsyncData<Nothing>
    data class Err(val err: Throwable) : AsyncData<Nothing>
    data class Ok<out T>(val data: T) : AsyncData<T>
}

fun <T : Any> asyncDataMutFlow() = MutableStateFlow<AsyncData<T>>(None)

fun <T : Any> T.asyncDataOk() = AsyncData.Ok<T>(this)
fun <T : Throwable> T.asyncDataErr() = AsyncData.Err(this)
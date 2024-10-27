package ru.saime.nice_pea_chat.common


sealed interface Status<out T> {
    object None : Status<Nothing>
    object Loading : Status<Nothing>
    data class Err(val err: Throwable) : Status<Nothing>
    data class Data<out T>(val data: T) : Status<T>
}
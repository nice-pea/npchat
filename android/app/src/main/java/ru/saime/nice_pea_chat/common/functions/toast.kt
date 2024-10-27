package ru.saime.nice_pea_chat.common.functions

import android.content.Context
import android.widget.Toast
import androidx.compose.runtime.Composable
import androidx.compose.ui.platform.LocalContext

enum class ToastDuration(val value: Int) {
    SHORT(0),
    LONG(1),
}

@Composable
fun Toast(
    text: String,
    duration: ToastDuration = ToastDuration.SHORT
) {
    Toast.makeText(LocalContext.current, text, duration.value).show()
}

fun toast(
    text: String,
    context: Context,
    duration: ToastDuration = ToastDuration.SHORT
) {
    Toast.makeText(context, text, duration.value).show()
}
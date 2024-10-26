package ru.saime.nice_pea_chat.data.store

import android.content.Context
import androidx.core.content.edit


class AuthenticationStore(context: Context) {
    private val name = "common"
    private val sp = context.getSharedPreferences(name, Context.MODE_PRIVATE)

    private val tokenKey = "token"
    var token: String
        get() = sp.getString(tokenKey, "").orEmpty()
        set(value) {
            sp.edit { putString(tokenKey, value) }
        }

    private val keyKey = "key"
    var key: String
        get() = sp.getString(keyKey, "").orEmpty()
        set(value) {
            sp.edit { putString(keyKey, value) }
        }
}
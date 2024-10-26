package ru.saime.nice_pea_chat.data.store

import android.content.Context
import androidx.core.content.edit


class NpcClientStore(context: Context) {
    private val name = "npcClientStore"
    private val sp = context.getSharedPreferences(name, Context.MODE_PRIVATE)

    private val hostKey = "host"
    var host: String
        get() = sp.getString(hostKey, "").orEmpty()
        set(value) {
            sp.edit { putString(hostKey, value) }
        }

    private val portKey = "port"
    var port: Int
        get() = sp.getInt(portKey, 0)
        set(value) {
            sp.edit { putInt(portKey, value) }
        }

//    private val urlKey = "url"
//    var url: String
//        get() = sp.getString(urlKey, "").orEmpty()
//        set(value) {
//            sp.edit { putString(urlKey, value) }
//        }

    fun urlIsNotEmpty() = host == "" || port == 0
}
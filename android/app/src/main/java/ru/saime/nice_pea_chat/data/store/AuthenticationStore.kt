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

    private val profileIdKey = "profileId"
    private val profileUsernameKey = "profileUsername"
    fun profile(): Result<Profile> {
        return runCatching {
            profile ?: error("no saved profile")
        }
    }

    var profile: Profile?
        get() {
            val id = sp.getInt(profileIdKey, 0)
            val username = sp.getString(profileUsernameKey, null).orEmpty()
            if (id == 0 || username == "") {
                return null
            }
            return Profile(id = id, username = username)
        }
        set(value) {
            sp.edit {
                putInt(profileIdKey, value?.id ?: 0)
                putString(profileUsernameKey, value?.username)
            }
        }

    fun clear() = sp.edit { clear() }
}

data class Profile(
    val id: Int,
    val username: String,
)
package ru.saime.nice_pea_chat.data.repositories

import ru.saime.nice_pea_chat.data.api.AuthenticationApi
import ru.saime.nice_pea_chat.data.api.AuthnResult
import ru.saime.nice_pea_chat.data.api.LoginResult
import ru.saime.nice_pea_chat.data.store.NpcClientStore

class AuthenticationRepository(
    private val api: AuthenticationApi,
    private val npcStore: NpcClientStore,
) {
    suspend fun authn(token: String, server: String = npcStore.host): Result<AuthnResult> {
        return api.authn(
            server = server,
            token = token,
        )
    }

    suspend fun login(key: String, server: String = npcStore.host): Result<LoginResult> {
        return api.login(
            server = server,
            key = key,
        )
    }
}
package ru.saime.nice_pea_chat.data.repositories

import ru.saime.nice_pea_chat.data.api.NpcClientApi
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import ru.saime.nice_pea_chat.network.npcBaseUrl

class NpcClient(
    private val api: NpcClientApi,
    private val npcStore: NpcClientStore,
) {
    suspend fun healthCheck(url: String = ""): Result<Unit> {
        return api.health(url = npcBaseUrl(npcStore, url))
    }
}
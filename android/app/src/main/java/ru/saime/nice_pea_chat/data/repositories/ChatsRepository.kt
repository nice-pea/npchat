package ru.saime.nice_pea_chat.data.repositories

import com.skydoves.retrofit.adapters.result.mapSuspend
import ru.saime.nice_pea_chat.data.api.ChatsApi
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.model.Model

class ChatsRepository(
    private val api: ChatsApi,
    private val authnStore: AuthenticationStore,
) {
    suspend fun chats(): Result<List<Model.Chat>> {
        return api.chats(
            unreadForUserID = authnStore.profile?.id
        ).mapSuspend { it.chats }
    }

    suspend fun chat(id: Int): Result<Model.Chat> {
        return api.chats(
            unreadForUserID = authnStore.profile?.id,
            ids = listOf(id),
        ).mapSuspend {
            it.chats.firstOrNull()
                ?: error("not found")
        }
    }
}

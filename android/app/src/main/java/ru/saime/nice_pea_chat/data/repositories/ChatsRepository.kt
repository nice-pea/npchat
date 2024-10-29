package ru.saime.nice_pea_chat.data.repositories

import ru.saime.nice_pea_chat.data.api.ChatsApi
import ru.saime.nice_pea_chat.data.api.ChatsResponse

class ChatsRepository(
    private val api: ChatsApi,
) {
    suspend fun chats(): Result<List<ChatsResponse.Chat>> {
        return api.chats()
    }
}

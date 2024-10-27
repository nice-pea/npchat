package ru.saime.nice_pea_chat.data.repositories

import com.skydoves.retrofit.adapters.result.mapSuspend
import ru.saime.nice_pea_chat.data.api.ChatsApi
import ru.saime.nice_pea_chat.data.api.ChatsResponse
import java.time.OffsetDateTime

class ChatsRepository(
    private val api: ChatsApi,
) {
    suspend fun chats(): Result<List<Chat>> {
        return api.chats().mapSuspend { chats ->
            chats.map(::Chat)
        }
    }
}

data class Chat(
    val id: Int,
    val name: String,
    val createdAt: OffsetDateTime,
    val creatorId: Int,
) {
    constructor(apiChat: ChatsResponse.Chat) : this(
        id = apiChat.id,
        name = apiChat.name,
        createdAt = apiChat.createdAt,
        creatorId = apiChat.creatorId
    )
}
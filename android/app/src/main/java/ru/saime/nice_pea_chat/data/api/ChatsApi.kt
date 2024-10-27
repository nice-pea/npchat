package ru.saime.nice_pea_chat.data.api

import retrofit2.http.GET
import java.time.OffsetDateTime

interface ChatsApi {
    @GET("/chats")
    suspend fun chats(
    ): Result<List<ChatsResponse.Chat>>
}

object ChatsResponse {
    data class Chat(
        val id: Int,
        val name: String,
        val createdAt: OffsetDateTime,
        val creatorId: Int,
    )
}
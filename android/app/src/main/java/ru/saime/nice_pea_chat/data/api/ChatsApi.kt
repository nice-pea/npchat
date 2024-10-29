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
        val id: Long,
        val name: String,
        val createdAt: OffsetDateTime,
        val creatorId: Long,

        val creator: Creator?,
        val lastMessage: Message?,
    )

    data class Creator(
        val id: Int,
        val username: String,
        val createdAt: OffsetDateTime
    )

    data class Message(
        val id: Long,
        val chatId: Long,
        val text: String,
        val authorId: Long?,
        val replyToId: Long?,
        val editedAt: OffsetDateTime?,
        val removedAt: OffsetDateTime?,
        val createdAt: OffsetDateTime,

        val author: Author?,
        val replyTo: Reply?,
    )

    data class Author(
        val id: Int,
        val username: String,
        val createdAt: OffsetDateTime
    )

    data class Reply(
        val id: Long,
        val chatId: Long,
        val text: String,
        val authorId: Long?,
        val replyToId: Long?,
        val editedAt: OffsetDateTime?,
        val removedAt: OffsetDateTime?,
        val createdAt: OffsetDateTime,

        val author: Author?,
    )

}
package ru.saime.nice_pea_chat.model

import java.time.OffsetDateTime


object Model {
    data class User(
        val id: Int,
        val username: String,
        val createdAt: OffsetDateTime
    )

    data class Session(
        val id: Int,
        val userId: Int,
        val token: String,
        val createdAt: OffsetDateTime,
        val expiresAt: OffsetDateTime
    )

    data class Chat(
        val id: Int,
        val name: String,
        val createdAt: OffsetDateTime,
        val creatorId: Int,

        val creator: User?,
        val lastMessage: Message?,
        val unreadMessagesCount: Int? = 0,
    )

    data class Message(
        val id: Int,
        val chatId: Int,
        val text: String,
        val authorId: Int?,
        val replyToId: Int?,
        val editedAt: OffsetDateTime?,
        val removedAt: OffsetDateTime?,
        val createdAt: OffsetDateTime,

        val author: User?,
        val replyTo: Reply?,
    )

    data class Reply(
        val id: Int,
        val chatId: Int,
        val text: String,
        val authorId: Int?,
        val replyToId: Int?,
        val editedAt: OffsetDateTime?,
        val removedAt: OffsetDateTime?,
        val createdAt: OffsetDateTime,

        val author: User?,
    )
}
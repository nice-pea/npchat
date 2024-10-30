package ru.saime.nice_pea_chat.data.api.model

import ru.saime.nice_pea_chat.model.Model

object ApiModel {
    data class Authn(
        val user: Model.User,
        val session: Model.Session
    )

    data class Login(
        val user: Model.User,
        val session: Model.Session
    )

    data class Chats(
        val chats: List<Model.Chat>,
    )
}
package ru.saime.nice_pea_chat.data.api

import retrofit2.http.GET
import retrofit2.http.Query
import ru.saime.nice_pea_chat.model.Model

interface MessagesApi {
    @GET("/messages")
    suspend fun messages(
        @Query("ids") ids: List<Int> = emptyList(),
        @Query("chat_ids") chatIDs: List<Int> = emptyList(),
        @Query("author_ids") authorIDs: List<Int> = emptyList(),
        @Query("reply_to_ids") replyToIDs: List<Int> = emptyList(),
        @Query("around_id") aroundID: Int? = null,
        @Query("before_id") beforeID: Int? = null,
        @Query("after_id") afterID: Int? = null,
        @Query("limit") limit: Int? = null,
    ): Result<List<Model.Message>>
}
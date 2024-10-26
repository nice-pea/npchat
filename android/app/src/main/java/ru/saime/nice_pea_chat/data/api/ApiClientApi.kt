package ru.saime.nice_pea_chat.data.api

import retrofit2.Call
import retrofit2.http.GET
import retrofit2.http.Path


interface ApiClientApi {
    @GET("{server}/health")
    fun health(
        @Path("server", encoded = true) server: String
    ): Call<Unit>
}
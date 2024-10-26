package ru.saime.nice_pea_chat.data.repositories

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.saime.nice_pea_chat.data.api.ApiClientApi

class ApiClient(
    private val api: ApiClientApi,
) {
    suspend fun healthCheck(server: String): Result<Unit> {
        return withContext(Dispatchers.IO) {
            val response = api.health(server = server).execute()
            when {
                response.isSuccessful -> Result.success(Unit)
                else -> Result.failure(Error(response.message()))
            }
        }
    }

    fun updateBaseUrl(url: String) {
        TODO("retrofit don't support dynamic baseurl")
    }
}
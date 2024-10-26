package ru.saime.nice_pea_chat.data

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.saime.nice_pea_chat.data.network.api.AuthenticationApi
import ru.saime.nice_pea_chat.data.network.api.AuthnResult
import ru.saime.nice_pea_chat.data.network.api.LoginResult

class AuthenticationRepository(
    private val authnApi: AuthenticationApi,
) {
    suspend fun authn(server: String, token: String): Result<AuthnResult> {
        return withContext(Dispatchers.IO) {
            val response = authnApi.Authn(
                server = server,
                token = token,
            ).execute()
            if (response.isSuccessful) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Error(response.message()))
            }
        }
    }

    suspend fun login(server: String, key: String): Result<LoginResult> {
        return withContext(Dispatchers.IO) {
            val response = authnApi.Login(
                server = server,
                key = key,
            ).execute()
            if (response.isSuccessful) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Error(response.message()))
            }
        }
    }
}
package ru.saime.nice_pea_chat.screens.app.authentication

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.saime.nice_pea_chat.data.repositories.AuthenticationRepository
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.data.store.NpcClientStore


sealed interface CheckAuthnResult {
    object None : CheckAuthnResult
    object Successful : CheckAuthnResult
    object ErrNoSavedCreds : CheckAuthnResult
    data class Err(val msg: String) : CheckAuthnResult
}

sealed interface AuthenticationAction {
    object CheckAuthn : AuthenticationAction
    object CheckAuthnConsume : AuthenticationAction
}

class AuthenticationViewModel(
    private val store: AuthenticationStore,
    private val npcStore: NpcClientStore,
    private val repo: AuthenticationRepository,
) : ViewModel() {

    private val _checkAuthnResult = MutableStateFlow<CheckAuthnResult>(CheckAuthnResult.None)
    val checkAuthnResult = _checkAuthnResult.asStateFlow()

    fun action(action: AuthenticationAction) {
        when (action) {
            AuthenticationAction.CheckAuthn -> viewModelScope.launch { checkAuthn() }
            AuthenticationAction.CheckAuthnConsume -> _checkAuthnResult.update { CheckAuthnResult.None }
        }
    }

    private suspend fun checkAuthn() {
        if (listOf(npcStore.host, store.key).any(String::isBlank)) {
            _checkAuthnResult.value = CheckAuthnResult.ErrNoSavedCreds
            return
        }
        val res = repo.authn(server = npcStore.host, token = store.token)
        when {
            res.isSuccess -> {
                store.token = res.getOrThrow().session.token
                _checkAuthnResult.value = CheckAuthnResult.Successful
            }

            res.isFailure -> {
                res.exceptionOrNull()?.toString().orEmpty().ifBlank { "blankErr" }
                    .run(CheckAuthnResult::Err)
                    .let { _checkAuthnResult.value = it }
            }
        }
    }

}
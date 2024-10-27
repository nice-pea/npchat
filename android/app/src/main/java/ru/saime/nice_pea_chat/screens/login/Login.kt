package ru.saime.nice_pea_chat.screens.login

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.input.TextFieldState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import org.koin.androidx.compose.koinViewModel
import ru.saime.nice_pea_chat.data.repositories.AuthenticationRepository
import ru.saime.nice_pea_chat.data.repositories.NpcClient
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.data.store.NpcClientStore
import ru.saime.nice_pea_chat.ui.components.Button
import ru.saime.nice_pea_chat.ui.components.Input
import ru.saime.nice_pea_chat.ui.functions.ToastDuration
import ru.saime.nice_pea_chat.ui.functions.toast


@Preview(
    backgroundColor = 0xFF000000,
    showBackground = true,
)
@Composable
private fun PreviewLoginScreen() {
    LoginScreen(rememberNavController())
}

const val RouteLogin = "Login"

@Composable
fun LoginScreen(
    navController: NavController,
) {
    val loginVM = koinViewModel<LoginViewModel>()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(20.dp),
        verticalArrangement = Arrangement.Center,
    ) {
        Input(
            title = "Server",
            placeholder = "http://example.com",
            textFieldState = loginVM.serverFieldState
        )
        Button(
            onClick = { loginVM.action(LoginAction.CheckConn) },
            text = "Check connection",
        )
        Input(
            title = "Key",
            placeholder = "Enter key for access to server",
            textFieldState = loginVM.keyFieldState
        )
        Button(
            onClick = { loginVM.action(LoginAction.Enter) },
            text = "Enter",
        )
    }

    CheckConnResultEffect(loginVM)
    EnterResultEffect(loginVM, navController)
}

@Composable
private fun CheckConnResultEffect(
    loginVM: LoginViewModel,
) {
    val ctx = LocalContext.current
    val result = loginVM.checkConnResult.collectAsState().value
    LaunchedEffect(result) {
        when (result) {
            is CheckConnResult.Err -> toast(result.msg, ctx, ToastDuration.LONG)
            CheckConnResult.Successful -> toast("connection established", ctx)
            CheckConnResult.None -> {}
        }
        loginVM.action(LoginAction.CheckConnConsume)
    }
}

@Composable
private fun EnterResultEffect(
    loginVM: LoginViewModel,
    navController: NavController,
) {
    val ctx = LocalContext.current
    val result = loginVM.enterResult.collectAsState().value
    LaunchedEffect(result) {
        when (result) {
            is EnterResult.Err -> toast(result.msg, ctx, ToastDuration.LONG)
            EnterResult.Successful -> navController.navigate("todo")
            EnterResult.None -> {}
        }
        loginVM.action(LoginAction.CheckConnConsume)
    }
}


sealed interface CheckConnResult {
    object None : CheckConnResult
    object Successful : CheckConnResult
    data class Err(val msg: String) : CheckConnResult
}

sealed interface EnterResult {
    object None : EnterResult
    object Successful : EnterResult
    data class Err(val msg: String) : EnterResult
}

sealed interface LoginAction {
    object CheckConn : LoginAction
    object CheckConnConsume : LoginAction
    object Enter : LoginAction
    object EnterConsume : LoginAction
}


class LoginViewModel(
    private val authnRepo: AuthenticationRepository,
    private val authnStore: AuthenticationStore,
    private val npcStore: NpcClientStore,
    private val apiClient: NpcClient
) : ViewModel() {
    val serverFieldState = TextFieldState("")
    val keyFieldState = TextFieldState("")

    private val _checkConnResult = MutableStateFlow<CheckConnResult>(CheckConnResult.None)
    val checkConnResult = _checkConnResult.asStateFlow()

    private val _enterResult = MutableStateFlow<EnterResult>(EnterResult.None)
    val enterResult = _enterResult.asStateFlow()

    fun action(action: LoginAction) {
        when (action) {
            LoginAction.CheckConn -> viewModelScope.launch { checkConn() }
            LoginAction.CheckConnConsume -> _checkConnResult.update { CheckConnResult.None }
            LoginAction.Enter -> viewModelScope.launch { enter() }
            LoginAction.EnterConsume -> _enterResult.update { EnterResult.None }
        }
    }

    private suspend fun checkConn() {
        val server = serverFieldState.text.toString()
        val res = apiClient.healthCheck(server)
        when {
            res.isSuccess -> _checkConnResult.update { CheckConnResult.Successful }
            res.isFailure -> {
                res.exceptionOrNull()?.toString().orEmpty().ifEmpty { "emptyErr" }
                    .run(CheckConnResult::Err)
                    .let { _checkConnResult.value = it }
            }
        }
    }

    private suspend fun enter() {
        val server = serverFieldState.text.toString()
        val key = keyFieldState.text.toString()
        val res = authnRepo.login(key = key, server = server)
        when {
            res.isSuccess -> {
                authnStore.token = res.getOrThrow().session.token
                npcStore.baseUrl = server
                authnStore.key = key
                _enterResult.update { EnterResult.Successful }
            }

            res.isFailure -> {
                res.exceptionOrNull()?.toString().orEmpty().ifEmpty { "emptyErr" }
                    .run(EnterResult::Err)
                    .let { _enterResult.value = it }
            }
        }
    }

}
package ru.saime.nice_pea_chat.screens.chats

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
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
import ru.saime.nice_pea_chat.data.repositories.Chat
import ru.saime.nice_pea_chat.data.repositories.ChatsRepository
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.screens.login.LoginScreen
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.Font

@Preview(
    backgroundColor = 0xFF000000,
    showBackground = true,
)
@Composable
private fun PreviewChatsScreen() {
    LoginScreen(rememberNavController())
}

const val RouteChats = "Chats"

@Composable
fun ChatsScreen(
    navController: NavController,
) {
    val chatsVM = koinViewModel<ChatsViewModel>()
    val uiState = chatsVM.uiState.collectAsState().value
    when (uiState) {
        is ChatsUiState.Content -> Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Black),
        ) {
            uiState.chats.forEach { chat ->
                Gap(4.dp)
                HorizontalDivider()
                Text(chat.id.toString(), style = Font.White12W500)
                Text(chat.name, style = Font.White12W500)
                Text(chat.creatorId.toString(), style = Font.White12W500)
                Text(chat.createdAt.toString(), style = Font.White12W500)
                Gap(2.dp)
            }
        }

        is ChatsUiState.Err -> Text(uiState.msg, style = Font.White12W500)
        ChatsUiState.Loading -> CircularProgressIndicator()
    }
}

sealed interface ChatsUiState {
    object Loading : ChatsUiState
    data class Err(val msg: String) : ChatsUiState
    data class Content(
        val chats: List<Chat>,
    ) : ChatsUiState
}

class ChatsViewModel(
    private val repo: ChatsRepository,
    private val store: AuthenticationStore,
) : ViewModel() {
    private val _uiState = MutableStateFlow<ChatsUiState>(ChatsUiState.Loading)
    val uiState = _uiState.asStateFlow()

    init {
        viewModelScope.launch {
            loadChats()
        }
    }

    private suspend fun loadChats() {
        repo.chats()
            .onSuccess { res ->
                _uiState.update {
                    ChatsUiState.Content(
                        chats = res
                    )
                }
            }
            .onFailure { res ->
                _uiState.update { ChatsUiState.Err(res.message.orEmpty()) }
            }
    }
}
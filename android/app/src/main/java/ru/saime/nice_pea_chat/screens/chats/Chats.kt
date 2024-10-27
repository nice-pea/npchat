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
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import org.koin.androidx.compose.koinViewModel
import ru.saime.nice_pea_chat.common.Status
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
    val chats = chatsVM.uiState.chats.collectAsState().value
    when (chats) {
        is Status.Data -> Column(
            modifier = Modifier
                .fillMaxSize()
                .background(Black),
        ) {
            chats.data.forEach { chat ->
                Gap(4.dp)
                HorizontalDivider()
                Text(chat.id.toString(), style = Font.White12W500)
                Text(chat.name, style = Font.White12W500)
                Text(chat.creatorId.toString(), style = Font.White12W500)
                Text(chat.createdAt.toString(), style = Font.White12W500)
                Gap(2.dp)
            }
        }

        is Status.Err -> Text(chats.err.toString(), style = Font.White12W500)
        Status.Loading -> CircularProgressIndicator()
        Status.None -> {}
    }
}


class ChatsUiState(
    val chats: StateFlow<Status<List<Chat>>>,
)

class ChatsViewModel(
    private val repo: ChatsRepository,
    private val store: AuthenticationStore,
) : ViewModel() {
    private val chats = MutableStateFlow<Status<List<Chat>>>(Status.None)
    val uiState = ChatsUiState(chats = chats)

    init {
        viewModelScope.launch {
            loadChats()
        }
    }

    private suspend fun loadChats() {
        repo.chats()
            .onSuccess { res ->
                chats.value = Status.Data(res)
            }
            .onFailure { err ->
                chats.value = Status.Err(err)
            }
    }
}
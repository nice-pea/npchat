package ru.saime.nice_pea_chat.screens.chat.messages

import androidx.compose.runtime.Composable
import androidx.navigation.NavController
import org.koin.androidx.compose.koinViewModel
import org.koin.core.parameter.parametersOf
import ru.saime.nice_pea_chat.screens.chats.ChatsViewModel

data class RouteMessages(
    val chatID: Int,
)

@Composable
fun MessagesScreen(
    navController: NavController,
    chatID: Int,
) {
    koinViewModel<ChatsViewModel> { parametersOf(chatID) }


}

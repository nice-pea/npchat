package ru.saime.nice_pea_chat.screens.chats

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Add
import androidx.compose.material.icons.outlined.Person
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.tooling.preview.Preview
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.NavController
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import org.koin.androidx.compose.koinViewModel
import ru.saime.nice_pea_chat.common.AsyncData
import ru.saime.nice_pea_chat.data.repositories.ChatsRepository
import ru.saime.nice_pea_chat.data.store.AuthenticationStore
import ru.saime.nice_pea_chat.model.Model
import ru.saime.nice_pea_chat.screens.chat.messages.RouteMessages
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.components.Progress
import ru.saime.nice_pea_chat.ui.theme.Dp2
import ru.saime.nice_pea_chat.ui.theme.Dp20
import ru.saime.nice_pea_chat.ui.theme.Dp4
import ru.saime.nice_pea_chat.ui.theme.Dp8
import ru.saime.nice_pea_chat.ui.theme.Font
import ru.saime.nice_pea_chat.ui.theme.RoundMin
import ru.saime.nice_pea_chat.ui.theme.White
import java.time.OffsetDateTime

const val RouteChats = "Chats"
private const val AppBarTitle = "Chats"

@Composable
fun ChatsScreen(
    navController: NavController,
) {
    val vm = koinViewModel<ChatsViewModel>()
    val chats = vm.uiState.chats.collectAsState().value
    val selfID = vm.uiState.selfID.collectAsState().value
    Column {
        AppBar(
            clickPlus = { },
            clickProfile = { },
        )
        when (chats) {
            is AsyncData.Ok -> ChatList(
                chats = chats.data,
                selfID = selfID,
                clickChat = {
                    navController.navigate(RouteMessages(chatID = it.id))
                }
            )

            is AsyncData.Err -> Text(chats.err.toString(), style = Font.White12W500)
            AsyncData.Loading -> Progress()
            AsyncData.None -> {}
        }
    }
}


@Composable
private fun ChatList(
    clickChat: (Model.Chat) -> Unit,
    chats: List<Model.Chat>,
    selfID: Int,
) {
    Column(
        verticalArrangement = Arrangement.spacedBy(Dp2),
    ) {
        chats.forEach { chat ->
            ChatCard(
                chat = chat,
                selfID = selfID,
                onClick = { clickChat(chat) }
            )
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun AppBar(
    clickPlus: () -> Unit,
    clickProfile: () -> Unit,
    modifier: Modifier = Modifier,
) {
    TopAppBar(
        modifier = modifier,
        title = { Text(AppBarTitle, style = Font.White20W400) },
        actions = {
            IconButton(onClick = clickPlus) {
                Icon(
                    imageVector = Icons.Outlined.Add,
                    tint = White,
                    contentDescription = "add chat"
                )
            }
            IconButton(onClick = clickProfile) {
                Icon(
                    imageVector = Icons.Outlined.Person,
                    tint = White,
                    contentDescription = "profile"
                )
            }
        }
    )
}


class ChatsUiState(
    val chats: StateFlow<AsyncData<List<Model.Chat>>>,
    val selfID: StateFlow<Int>,
)


class ChatsViewModel(
    private val repo: ChatsRepository,
    private val store: AuthenticationStore,
) : ViewModel() {
    private val chats = MutableStateFlow<AsyncData<List<Model.Chat>>>(AsyncData.None)
    private val selfID = MutableStateFlow<Int>(0)
    val uiState = ChatsUiState(chats = chats, selfID = selfID)

    init {
        viewModelScope.launch {
            loadChats()
        }
        viewModelScope.launch {
            selfID.value = store.profile?.id ?: 0
        }
    }

    private suspend fun loadChats() {
        repo.chats().onSuccess { res ->
            chats.value = AsyncData.Ok(res)
        }.onFailure { err ->
            chats.value = AsyncData.Err(err)
        }
    }
}


@Preview(
    backgroundColor = 0xFF000000,
    showBackground = true,
)
@Composable
private fun PreviewChatCard() {
    val authorID = 51342
    ChatCard(
        onClick = {},
        selfID = authorID, chat = Model.Chat(
            id = 123,
            name = "ChatName",
            createdAt = OffsetDateTime.now(),
            creatorId = 41,
            creator = null,
            lastMessage = Model.Message(
                id = 23,
                chatId = 123,
                text = "Message text",
                authorId = authorID,
                replyToId = 123,
                editedAt = null,
                removedAt = null,
                createdAt = OffsetDateTime.now(),
                author = Model.User(
                    id = authorID, username = "UserName", createdAt = OffsetDateTime.now()
                ),
                replyTo = null
            ),
            unreadMessagesCount = 23
        )
    )
}


@Composable
private fun ChatCard(
    modifier: Modifier = Modifier,
    onClick: () -> Unit,
    chat: Model.Chat,
    selfID: Int,
) {
    val noMessagesPlaceholder = "<No messages>"
    val lastMessage = chat.lastMessage
    val authorID = lastMessage?.authorId
    val unreadCount = chat.unreadMessagesCount
    val authorTextStyle = when (authorID) {
        selfID -> Font.BlueSky14W400
        else -> Font.White14W400
    }
    val messageTextStyle = when {
        lastMessage == null -> Font.DarkGraph14W500Italic
        lastMessage.removedAt != null -> Font.DarkGraph14W500Italic
        lastMessage.author == null -> Font.Pink14W500Italic
        else -> Font.BlueGraph14W500
    }

    Column(
        modifier = Modifier
            .clip(RoundMin)
            .clickable(onClick = onClick)
            .padding(horizontal = Dp20, vertical = Dp8)
            .then(modifier),
    ) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
        ) {
            Text(
                modifier = Modifier.weight(1f),
                text = chat.name,
                maxLines = 1,
                style = Font.White14W400,
            )
            if (unreadCount != null && unreadCount > 0) {
                Gap(Dp4)
                Text("+${unreadCount}", style = Font.BlueGraph12W400)
            }
        }
        if (lastMessage == null) Text(
            modifier = Modifier.padding(start = Dp8),
            text = noMessagesPlaceholder,
            style = messageTextStyle,
        )
        else Row(
            modifier = Modifier.padding(start = Dp8),
            verticalAlignment = Alignment.CenterVertically,
        ) {
            if (lastMessage.author != null) {
                Text(
                    text = lastMessage.author.username,
                    maxLines = 1,
                    style = authorTextStyle,
                )
            }
            Gap(Dp4)
            Text(
                modifier = Modifier.weight(1f),
                text = lastMessage.text,
                maxLines = 1,
                style = messageTextStyle,
            )
            Gap(Dp4)
            Text(
                text = lastMessage.createdAt.run { "$hour:$minute" },
                style = Font.DarkGraph12W400,
            )
        }
    }
}
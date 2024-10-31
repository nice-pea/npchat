package ru.saime.nice_pea_chat.screens.app

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.toRoute
import org.koin.core.Koin
import ru.saime.nice_pea_chat.screens.chat.messages.MessagesScreen
import ru.saime.nice_pea_chat.screens.chat.messages.RouteMessages
import ru.saime.nice_pea_chat.screens.chats.ChatsScreen
import ru.saime.nice_pea_chat.screens.chats.RouteChats
import ru.saime.nice_pea_chat.screens.login.LoginScreen
import ru.saime.nice_pea_chat.screens.login.RouteLogin
import ru.saime.nice_pea_chat.screens.splash.RouteSplash
import ru.saime.nice_pea_chat.screens.splash.SplashScreen
import ru.saime.nice_pea_chat.ui.theme.Black


@Composable
fun ComposeApp(koin: Koin) {
    val navController = rememberNavController()
    NavHost(
        modifier = Modifier
            .fillMaxSize()
            .background(Black),
        navController = navController,
        startDestination = RouteSplash
    ) {
        composable(RouteSplash) { SplashScreen(navController) }
        composable(RouteLogin) { LoginScreen(navController) }
        composable(RouteChats) { ChatsScreen(navController) }
        composable<RouteMessages> {
            val messages = it.toRoute<RouteMessages>()
            MessagesScreen(navController, messages.chatID)
        }
    }
}

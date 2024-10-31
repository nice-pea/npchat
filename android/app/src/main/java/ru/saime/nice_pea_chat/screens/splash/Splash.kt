package ru.saime.nice_pea_chat.screens.splash

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.navigation.NavController
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.delay
import org.koin.androidx.compose.koinViewModel
import ru.saime.nice_pea_chat.common.functions.toast
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationAction
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationViewModel
import ru.saime.nice_pea_chat.screens.app.authentication.CheckAuthnResult
import ru.saime.nice_pea_chat.screens.chats.RouteChats
import ru.saime.nice_pea_chat.screens.login.RouteLogin
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.components.Progress
import ru.saime.nice_pea_chat.ui.modifiers.fadeIn
import ru.saime.nice_pea_chat.ui.theme.Dp10
import ru.saime.nice_pea_chat.ui.theme.Font
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds
import kotlin.time.Duration.Companion.seconds


@Preview()
@Composable
private fun PreviewSplashScreen() {
    SplashScreen(
        navController = rememberNavController(),
        textFadeInDuration = 0.milliseconds,
        loaderFadeInDuration = 0.milliseconds,
    )
}

const val RouteSplash = "Splash"
private const val Title = "nice-pea-chat\n(NPC)"

@Composable
fun SplashScreen(
    navController: NavController,
    textFadeInDuration: Duration = 300.milliseconds,
    loaderFadeInDuration: Duration = 300.milliseconds,
) {
    Column(
        modifier = Modifier.fillMaxSize(),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            modifier = Modifier.fadeIn(textFadeInDuration),
            text = Title,
            style = Font.White16W400,
            textAlign = TextAlign.Center
        )
        Gap(Dp10)
        Progress(modifier = Modifier.fadeIn(loaderFadeInDuration))
    }

    val vm = koinViewModel<AuthenticationViewModel>()
    CheckAuthnResultEffect(navController, vm)
    LaunchedEffect(1) {
        vm.action(AuthenticationAction.CheckAuthn)
    }
}

@Composable
private fun CheckAuthnResultEffect(
    navController: NavController,
    authnVM: AuthenticationViewModel,
) {
    val ctx = LocalContext.current
    val checkAuthnResult = authnVM.checkAuthnResult.collectAsState().value
    LaunchedEffect(checkAuthnResult) {
        delay(.7.seconds)
        when (checkAuthnResult) {
            is CheckAuthnResult.Err -> toast(checkAuthnResult.msg, ctx)
            CheckAuthnResult.ErrNoSavedCreds -> navController.navigate(RouteLogin)
            CheckAuthnResult.Successful -> navController.navigate(RouteChats)
            CheckAuthnResult.None -> {}
        }
        authnVM.action(AuthenticationAction.CheckAuthnConsume)
    }
}
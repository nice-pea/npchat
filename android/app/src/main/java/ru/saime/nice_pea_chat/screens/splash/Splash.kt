package ru.saime.nice_pea_chat.screens.splash

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.delay
import org.koin.androidx.compose.koinViewModel
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationAction
import ru.saime.nice_pea_chat.screens.app.authentication.AuthenticationViewModel
import ru.saime.nice_pea_chat.screens.app.authentication.CheckAuthnResult
import ru.saime.nice_pea_chat.screens.login.RouteLogin
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.functions.toast
import ru.saime.nice_pea_chat.ui.modifiers.fadeIn
import ru.saime.nice_pea_chat.ui.theme.Font
import ru.saime.nice_pea_chat.ui.theme.White
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
            text = "nice-pea-chat\n(NPC)",
            style = Font.White16W400,
            textAlign = TextAlign.Center
        )
        Gap(10.dp)
        CircularProgressIndicator(
            modifier = Modifier.fadeIn(loaderFadeInDuration),
            color = White,
        )
    }

    val authnVM = koinViewModel<AuthenticationViewModel>()
    CheckAuthnResultEffect(navController, authnVM)
    LaunchedEffect(1) {
        authnVM.action(AuthenticationAction.CheckAuthn)
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
        when (checkAuthnResult) {
            is CheckAuthnResult.Err -> toast(checkAuthnResult.msg, ctx)
            CheckAuthnResult.ErrNoSavedCreds -> {
                delay(.7.seconds)
                navController.navigate(RouteLogin)
                toast("ErrNoSavedCreds", ctx)
            }

            CheckAuthnResult.Successful -> {
                delay(.7.seconds)
                navController.navigate(RouteLogin)
                toast("ErrNoSavedCreds", ctx)
            }

            CheckAuthnResult.None -> {}
        }
        authnVM.action(AuthenticationAction.CheckAuthnConsume)
    }
}
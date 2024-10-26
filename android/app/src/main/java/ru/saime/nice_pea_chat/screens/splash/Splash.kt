package ru.saime.nice_pea_chat.screens.splash

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.modifiers.fadeIn
import ru.saime.nice_pea_chat.ui.theme.Font
import ru.saime.nice_pea_chat.ui.theme.White
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds


@Preview()
@Composable
private fun PreviewSplashScreen() {
    SplashScreen(
        textFadeInDuration = 0.milliseconds,
        loaderFadeInDuration = 0.milliseconds,
        job = {},
    )
}


@Composable
fun SplashScreen(
    textFadeInDuration: Duration = 300.milliseconds,
    loaderFadeInDuration: Duration = 300.milliseconds,
    job: suspend () -> Unit,
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
    LaunchedEffect(1) {
        job()
    }
}

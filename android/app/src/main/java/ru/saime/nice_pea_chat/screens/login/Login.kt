package ru.saime.nice_pea_chat.screens.login

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.input.rememberTextFieldState
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import ru.saime.nice_pea_chat.ui.components.Button
import ru.saime.nice_pea_chat.ui.components.Input


@Preview(
    backgroundColor = 0xFF000000,
    showBackground = true,
)
@Composable
private fun PreviewLoginScreen() {
    LoginScreen()
}
const val RouteLogin = "Login"

@Composable
fun LoginScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(20.dp),
        verticalArrangement = Arrangement.Center,
    ) {
        Input(
            title = "Server",
            placeholder = "http://example.com",
            textFieldState = rememberTextFieldState()
        )
        Button(
            onClick = {},
            text = "Check connection",
        )
        Input(
            title = "Key",
            placeholder = "Enter key for access to server",
            textFieldState = rememberTextFieldState()
        )
        Button(
            onClick = {},
            text = "Enter",
        )
    }
}
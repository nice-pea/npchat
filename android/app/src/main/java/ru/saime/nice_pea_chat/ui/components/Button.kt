package ru.saime.nice_pea_chat.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.tooling.preview.Preview
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.Dp10
import ru.saime.nice_pea_chat.ui.theme.Dp2
import ru.saime.nice_pea_chat.ui.theme.Dp20
import ru.saime.nice_pea_chat.ui.theme.Dp8
import ru.saime.nice_pea_chat.ui.theme.Font
import ru.saime.nice_pea_chat.ui.theme.RoundMin


@Preview
@Composable
private fun PreviewButton() {
    Button(
        modifier = Modifier
            .background(Black)
            .padding(Dp20),
        text = "Confirm",
        helperText = "The number of chats that can be created is limited. Created chats cannot be deleted",
        onClick = {}
    )
}

@Composable
fun Button(
    modifier: Modifier = Modifier,
    onClick: () -> Unit,
    text: String,
    helperText: String = "",
) {
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundMin)
            .clickable(onClick = onClick)
            .padding(Dp8)
            .then(modifier)
    ) {
        Row {
            Text("->", style = Font.White16W400)
            Gap(Dp10)
            Text(text.ifEmpty { "<action>" }, style = Font.White16W400)
        }
        if (helperText != "") {
            Gap(Dp2)
            Text(helperText, style = Font.GrayCharcoal12W400)
        }
    }
}
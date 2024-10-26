package ru.saime.nice_pea_chat.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.Font


@Preview
@Composable
private fun PreviewButton() {
    Button(
        modifier = Modifier
            .background(Black)
            .padding(20.dp),
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
            .clip(RoundedCornerShape(3.dp))
            .clickable(onClick = onClick)
            .padding(8.dp)
            .then(modifier)
    ) {
        Row {
            Text("->", style = Font.White16W400)
            Gap(10.dp)
            Text(text.ifBlank { "<action>" }, style = Font.White16W400)
        }
        if (helperText.isNotBlank()) {
            Gap(2.dp)
            Text(helperText, style = Font.GrayCharcoal12W400)
        }
    }
}
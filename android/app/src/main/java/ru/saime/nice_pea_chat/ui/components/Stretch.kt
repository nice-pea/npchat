package ru.saime.nice_pea_chat.ui.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.RowScope
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier

@Composable
fun ColumnScope.Stretch(modifier: Modifier = Modifier) {
    Box(
        modifier = Modifier
            .weight(1f)
            .then(modifier)
    )
}

@Composable
fun RowScope.Stretch(modifier: Modifier = Modifier) {
    Box(
        modifier = Modifier
            .weight(1f)
            .then(modifier)
    )
}
package ru.saime.nice_pea_chat.ui.modifiers

import androidx.compose.animation.core.Animatable
import androidx.compose.animation.core.tween
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.composed
import androidx.compose.ui.draw.alpha
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds

fun Modifier.fadeIn(duration: Duration = 300.milliseconds): Modifier = composed {
    val animatedFloat = remember { Animatable(0f) }
    val durationMillis = duration.inWholeMilliseconds.toInt()
    LaunchedEffect(Unit) {
        animatedFloat.animateTo(1f, animationSpec = tween(durationMillis = durationMillis))
    }

    this.then(
        Modifier.alpha(animatedFloat.value)
    )
}
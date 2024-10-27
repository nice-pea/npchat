package ru.saime.nice_pea_chat.ui.theme

import android.app.Activity
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.LocalRippleConfiguration
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.RippleConfiguration
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat

private val DarkColorScheme = darkColorScheme(
    primary = Copper,
    onPrimary = Olive,
    secondary = GrayCharcoal,
    tertiary = Pink,

    surface = Black,
    onSurface = White,
    onSurfaceVariant = GrayCharcoal,
)

private val LightColorScheme = lightColorScheme(
    primary = Copper,
    onPrimary = Olive,
    secondary = GrayCharcoal,
    tertiary = Pink,

    surface = White,
    onSurface = Black,
    onSurfaceVariant = GrayCharcoal,
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun NicePeaChatTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) DarkColorScheme else LightColorScheme

    // Make status bar color as contrast
    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.background.toArgb()
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = !darkTheme
        }
    }
    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
    ) {
        CompositionLocalProvider(
            LocalRippleConfiguration provides rippleConfiguration,
            content = content
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
private val rippleConfiguration = RippleConfiguration(color = BlueGraph)

val cursorBrush = SolidColor(Copper)
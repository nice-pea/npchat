package ru.saime.nice_pea_chat

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin
import ru.saime.nice_pea_chat.di.appModule
import ru.saime.nice_pea_chat.screens.app.ComposeApp
import ru.saime.nice_pea_chat.ui.theme.NicePeaChatTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        installSplashScreen()

        val koinApp = startKoin {
            androidContext(this@MainActivity)
            modules(appModule)
        }

        setContent {
            NicePeaChatTheme {
                ComposeApp(koinApp.koin)
            }
        }
    }
}

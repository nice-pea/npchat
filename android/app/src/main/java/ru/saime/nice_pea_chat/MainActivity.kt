package ru.saime.nice_pea_chat

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import org.koin.android.ext.koin.androidContext
import org.koin.androidx.viewmodel.ext.android.viewModel
import org.koin.core.context.startKoin
import ru.saime.nice_pea_chat.di.appModule
import ru.saime.nice_pea_chat.ui.Gap
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.NicePeaChatTheme
import ru.saime.nice_pea_chat.ui.theme.White
import kotlin.time.Duration.Companion.seconds

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        startKoin {
            androidContext(this@MainActivity)
            modules(appModule)
        }

        val mainViewModel: MainViewModel by viewModel()

        installSplashScreen().setKeepOnScreenCondition {
            !mainViewModel.loaded.value
        }
//        enableEdgeToEdge()
        setContent {
            NicePeaChatTheme {
                ComposeApp()
            }
        }
    }
}


// ComposeApp.kt

enum class Route {
    Splash,
    List
}

@Composable
fun ComposeApp() {
    val navController = rememberNavController()
    NavHost(
        navController = navController,
        startDestination = Route.Splash.name
    ) {
        composable(Route.Splash.name) {
            SplashScreen(
                job = {
                    delay(2.seconds)
                },
                action = {
                    navController.navigate(Route.List.name)
                }
            )
        }
        composable(Route.List.name) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(Black),
                contentAlignment = Alignment.Center,
            ) {
                Column {
                    repeat(10) {
                        Text(text = "$it", color = White)
                        Gap(10.dp)
                    }
                }
            }
        }
    }
}

// Splash.kt

sealed interface SplashAction {
    data object OnJobFinal : SplashAction
}

@Composable
fun SplashScreen(
    job: suspend () -> Unit,
    action: (SplashAction) -> Unit,
) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Black),
        contentAlignment = Alignment.Center,
    ) {
        CircularProgressIndicator(
            color = White
        )
    }
    LaunchedEffect(1) {
        job()
        action(SplashAction.OnJobFinal)
    }
}



class MainViewModel() : ViewModel() {
    private val _loaded = MutableStateFlow(false)
    val loaded = _loaded.asStateFlow()
    init {
        viewModelScope.launch {
            delay(5.seconds)
            _loaded.value = true
        }
    }
}

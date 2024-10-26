package ru.saime.nice_pea_chat

import android.content.Context
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.core.content.edit
import androidx.core.splashscreen.SplashScreen.Companion.installSplashScreen
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.delay
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import org.koin.android.ext.koin.androidContext
import org.koin.core.Koin
import org.koin.core.context.startKoin
import retrofit2.Call
import retrofit2.http.GET
import retrofit2.http.Query
import ru.saime.nice_pea_chat.di.appModule
import ru.saime.nice_pea_chat.screens.login.LoginScreen
import ru.saime.nice_pea_chat.ui.components.Gap
import ru.saime.nice_pea_chat.ui.modifiers.fadeIn
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.Font
import ru.saime.nice_pea_chat.ui.theme.NicePeaChatTheme
import ru.saime.nice_pea_chat.ui.theme.White
import java.time.LocalDateTime
import kotlin.time.Duration
import kotlin.time.Duration.Companion.milliseconds
import kotlin.time.Duration.Companion.seconds

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        installSplashScreen()

        val koinApp = startKoin {
            androidContext(this@MainActivity)
            modules(appModule)
        }
//        enableEdgeToEdge()
        setContent {
            NicePeaChatTheme {
                ComposeApp(koinApp.koin)
            }
        }
    }
}

// ComposeApp.kt

enum class Route {
    Splash,
    Login
}

@Composable
fun ComposeApp(koin: Koin) {
    val navController = rememberNavController()
    NavHost(
        modifier = Modifier.background(MaterialTheme.colorScheme.surface),
        navController = navController,
        startDestination = Route.Splash.name
    ) {
        composable(Route.Splash.name) {
            SplashScreen(
                textFadeInDuration = 1.5.seconds,
                job = {
                    val authnRepo = koin.get<AuthenticationRepository>()
                    val authStore = koin.get<AuthnStore>()
                    authStore.SaveToken("f1d727b2-212e-47cf-a7a2-40ff581bc816")
                    delay(2.seconds)
                    when (val result = authnRepo.Authn(authStore.Token())) {
                        is AuthenticationRepository.CheckResult.Failed ->
                            Log.d("", result.error)

                        is AuthenticationRepository.CheckResult.Ok -> {
                            Log.d("", result.toString())
                            navController.navigate(Route.Login.name)
                        }
                    }

                },
            )
        }
        composable(Route.Login.name) {
            LoginScreen()
        }
    }
}

// Authentication.kt

data class User(
    val id: Int,
    val username: String,
    val createdAt: LocalDateTime // Используем LocalDateTime для представления времени
)

data class Session(
    val id: Int,
    val userId: Int,
    val token: String,
    val createdAt: LocalDateTime,
    val expiresAt: LocalDateTime
)

data class AuthnResult(
    val user: User,
    val session: Session
)

interface AuthenticationApi {
    @GET("/authn")
    fun Authn(@Query("token") token: String): Call<AuthnResult>

    @GET("/authn/login")
    fun Login(@Query("key") key: String): Call<AuthnResult>
}

class AuthenticationRepository(
    private val authnApi: AuthenticationApi,
) {
    sealed interface CheckResult {
        data class Ok(val data: AuthnResult) : CheckResult
        data class Failed(val error: String) : CheckResult
    }

    suspend fun Authn(token: String): CheckResult {
        return withContext(Dispatchers.IO) {
            val response = authnApi.Authn(token).execute()
            when {
//            response == null -> CheckResult.Failed("bull")
                response.isSuccessful -> CheckResult.Ok(response.body()!!)
                else -> CheckResult.Failed(response.message())
//            response.isFailure -> CheckResult.Failed(response.exceptionOrNull()!!.message.orEmpty())
//            response.getOrThrow().isSuccessful -> CheckResult.Ok
//            else -> CheckResult.Failed(response.getOrThrow().message())
            }
        }
    }
}


//sealed interface AuthenticationEvent {
//    //
//    data object OnAuthnFailed : AuthenticationEvent
//}
//
//class AuthenticationViewModel(
//    val AuthnRepo: AuthnRepository,
//) : ViewModel() {
//
//}


// AuthnStore.kt

//interface AuthnStore {
//    fun Token(): String
//    fun SaveToken(toke: String)
//}

//class AuthnStoreSharedPrefs {
class AuthnStore(context: Context) {

    private val sp = context.getSharedPreferences("common", Context.MODE_PRIVATE)

    fun Token(): String {
        return sp.getString("token", "").orEmpty()
    }

    fun SaveToken(token: String) {
        sp.edit {
            this.putString("token", token)
        }
    }
}

// Splash.kt

@Composable
fun SplashScreen(
    textFadeInDuration: Duration = 300.milliseconds,
    loaderFadeInDuration: Duration = 300.milliseconds,
    job: suspend () -> Unit,
) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(Black),
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


class MainViewModel : ViewModel() {
    private val _loaded = MutableStateFlow(false)
    val loaded = _loaded.asStateFlow()

    init {
        viewModelScope.launch {
            delay(5.seconds)
            _loaded.value = true
        }
    }
}

package ru.saime.nice_pea_chat

import android.content.Context
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.SystemBarStyle
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.background
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
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
import ru.saime.nice_pea_chat.screens.splash.SplashScreen
import ru.saime.nice_pea_chat.ui.theme.Black
import ru.saime.nice_pea_chat.ui.theme.NicePeaChatTheme
import java.time.LocalDateTime
import kotlin.time.Duration.Companion.seconds

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        installSplashScreen()

        val koinApp = startKoin {
            androidContext(this@MainActivity)
            modules(appModule)
        }
        enableEdgeToEdge(
            statusBarStyle = SystemBarStyle.dark(Color.Transparent.toArgb())
        )
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
        modifier = Modifier.background(Black),
        navController = navController,
        startDestination = Route.Splash.name
    ) {
        composable(Route.Splash.name) {
            SplashScreen(
                job = {
                    val authStore = koin.get<AuthnStore>()
                    if (authStore.Token().isBlank()) {
                        navController.navigate(Route.Login.name)
                        return@SplashScreen
                    }
                    val authnRepo = koin.get<AuthenticationRepository>()
                    delay(.7.seconds)
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

package registerHandler

import (
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

var JWTSecret = os.Getenv("OAUTH_JWT_SECRET")

func GoogleRegistration(router http2.Router, discovery adapter.ServiceDiscovery) {
	router.HandleFunc(
		"GET /oauth/google/registration",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			oauthState, err := GenerateStateOauthJWT(JWTSecret)
			if err != nil {
				return nil, err
			}
			//config := newGoogleOAuth2RegistrationConfig(discovery)

			return http2.Redirect{
				URL:  context.Adapters().OAuthGoogle.AuthCodeURL(oauthState),
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

func GoogleRegistrationCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/google/registration/callback",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			ctx := context.Request().Context()

			state := http2.FormStr(context, "state")
			if state == "" {
				return nil, errors.New("state is empty")
			}
			if err := ValidateStateOauthJWT(state, JWTSecret); err != nil {
				return nil, err
			}
			input := service.GoogleRegistrationInput{
				Code: http2.FormStr(context, "code"),
			}
			out, err := context.Services().OAuth().GoogleRegistration(input)
			if err != nil {
				return nil, err
			}

			data, err := gc.GoogleUseCase.GetUserDataFromGoogle(googleOauthConfig, r.FormValue("code"), oauthGoogleUrlAPI)
			if err != nil {
				log.Error(err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			code := http2.FormStr(context, "code")
			accessToken, refreshToken, err := gc.GoogleUseCase.GoogleLogin(ctx, data, gc.Env)
			if err != nil {
				log.Error(err)
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
		},
	)
}

//func newGoogleOAuth2RegistrationConfig(discovery adapter.ServiceDiscovery) *oauth2.Config {
//	return &oauth2.Config{
//		ClientID:     os.Getenv("GOOGLE_KEY"),
//		ClientSecret: os.Getenv("GOOGLE_SECRET"),
//		Endpoint:     googleOAuth.Endpoint,
//		RedirectURL:  discovery.NpcApiPubUrl() + "/oauth/google/registration/callback",
//		Scopes: []string{
//			"https://www.googleapis.com/auth/userinfo.email",
//			"https://www.googleapis.com/auth/userinfo.profile",
//		},
//	}
//}

func GenerateStateOauthJWT(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // State живёт 10 минут
	})
	return token.SignedString([]byte(secret))
}

func ValidateStateOauthJWT(stateToken, secret string) error {
	token, err := jwt.Parse(stateToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is not valid")
	}

	return nil
}

//func generateStateOauthCookie(w http.ResponseWriter) string {
//	expiration := time.Now().Add(365 * 24 * time.Hour)
//
//	b := make([]byte, 16)
//	rand.Read(b)
//	state := base64.URLEncoding.EncodeToString(b)
//	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
//	http.SetCookie(w, &cookie)
//
//	return state
//}
//func registerOAuthCallbacks(discovery adapter.ServiceDiscovery) {
//	config := newGoogleOAuth2RegistrationConfig(discovery)
//	relyingParty, err := rp.NewRelyingPartyOAuth(config)
//	if err != nil {
//		return
//	}
//
//	token := cli.CodeFlow[*oidc.IDTokenClaims](
//		context.Background(),
//		relyingParty,
//		"/oauth/google/callback",
//		"8080",
//		func() string {
//			return uuid.NewString()
//		},
//	)
//	token.IDTokenClaims
//
//	client, _ := googleOAuth.DefaultClient(context.Background(),
//		"https://www.googleapis.com/auth/userinfo.email",
//		"https://www.googleapis.com/auth/userinfo.profile",
//	)
//	client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
//}

//func OAuthAddMethod(router http2.Router) {
//
//}

//func OAuthCallback(router http2.Router, discovery adapter.ServiceDiscovery) {
//	registerOAuthCallbacks(discovery)
//	router.HandleFunc("GET /oauth/{provider}/callback", nil, func(context http2.Context) (any, error) {
//
//		return nil, nil
//	})
//
//}
//func OAuthCallbackHttp() rp.CodeExchangeCallback {
//	return func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
//
//	}
//}

//type OAuthContext struct {
//	http2.Context
//	Tokens *oidc.Tokens
//	State  string
//}
//
//func OAuthExchangerToHandleFunc(hf func(OAuthContext) (any, error)) func(http.ResponseWriter, *http.Request, *oidc.Tokens, string, rp.RelyingParty) {
//	return func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
//		ctx := OAuthContext{
//			Context: nil,
//			Tokens:  nil,
//			State:   "",
//		}
//	}
//}

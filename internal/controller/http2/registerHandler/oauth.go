package registerHandler

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

func OAuthInitRegistration(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.OAuthInitRegistrationInput{
				Provider: http2.PathStr(context, "provider"),
			}
			out, err := context.Services().OAuth().InitRegistration(input)
			if err != nil {
				return nil, err
			}

			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return nil, err
			}

			return http2.Redirect{
				URL:  out.RedirectURL,
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

func OAuthCompleteRegistrationCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration/callback",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			if err := validateOAuthCookie(context); err != nil {
				return nil, err
			}

			input := service.OAuthCompeteRegistrationInput{
				UserCode: http2.FormStr(context, "code"),
				Provider: http2.PathStr(context, "provider"),
			}

			return context.Services().OAuth().CompeteRegistration(input)
		},
	)
}

func OAuthInitLogin(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/login",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.OAuthInitLoginInput{
				Provider: http2.PathStr(context, "provider"),
			}
			out, err := context.Services().OAuth().InitLogin(input)
			if err != nil {
				return nil, err
			}

			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return nil, err
			}

			return http2.Redirect{
				URL:  out.RedirectURL,
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

func OAuthCompleteLoginCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/login/callback",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			if err := validateOAuthCookie(context); err != nil {
				return nil, err
			}

			input := service.OAuthCompleteLoginInput{
				UserCode:  http2.FormStr(context, "code"),
				InitState: http2.FormStr(context, "state"),
				Provider:  http2.PathStr(context, "provider"),
			}

			return context.Services().OAuth().CompleteLogin(input)
		},
	)
}

const oauthCookieName = "oauthState"

func setOAuthCookie(context http2.Context, redirectURL string) error {
	parsedUrl, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}

	http.SetCookie(context.Writer(), &http.Cookie{
		Name:    oauthCookieName,
		Value:   parsedUrl.Query().Get("state"),
		Expires: time.Now().Add(time.Hour),
	})

	return nil
}

var errWrongState = errors.New("неправильный state")

func validateOAuthCookie(context http2.Context) error {
	oauthState, err := context.Request().Cookie(oauthCookieName) // игнорировать ErrNoCookie
	if errors.Is(err, http.ErrNoCookie) {
		return errWrongState
	} else if err != nil {
		return err
	}

	if http2.FormStr(context, "state") != oauthState.Value {
		return errWrongState
	}

	return nil
}

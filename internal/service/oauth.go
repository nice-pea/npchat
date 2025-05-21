package service

import (
	"errors"
	"fmt"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Google    adapter.OAuthGoogle
	OAuthRepo domain.OAuthRepository
}

type OAuthRegistrationCallbackInput struct {
	Provider string
}

type GoogleRegistrationInput struct {
	UserCode  string
	InitState string
}

//type GoogleRegistrationOut struct {
//	User    domain.User
//	Session domain.Session
//}

// Validate валидирует значение отдельно каждого параметры
func (in GoogleRegistrationInput) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.InitState == "" {
		return ErrInvalidInitState
	}

	return nil
}

// GoogleRegistration
// Подсмотрено в: https://github.com/oguzhantasimaz/Go-Clean-Architecture-Template/blob/main/api/controller/google.go
func (o *OAuth) GoogleRegistration(in GoogleRegistrationInput) (domain.User, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.User{}, err
	}

	// Проверить InitState
	links, err := o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{ID: in.InitState})
	if err != nil {
		return domain.User{}, err
	}
	if len(links) != 1 {
		return domain.User{}, ErrWrongInitState
	}

	token, err := o.Google.Exchange(in.UserCode)
	if err != nil {
		return domain.User{}, errors.Join(ErrWrongUserCode, err)
	}
	googleUser, err := o.Google.User(token)
	if err != nil {
		return domain.User{}, err
	}

	fmt.Println(googleUser)

	return domain.User{}, nil
}

type GoogleRegistrationInitOut struct {
	RedirectURL string
}

//var JWTSecret = os.Getenv("OAUTH_JWT_SECRET")

func (o *OAuth) GoogleRegistrationInit() (GoogleRegistrationInitOut, error) {
	//oauthState, err := generateStateOauthJWT(JWTSecret)
	//if err != nil {
	//	return GoogleRegistrationInitOut{}, err
	//}

	return GoogleRegistrationInitOut{
		RedirectURL: o.Google.AuthCodeURL(""),
	}, nil
}

//func generateStateOauthJWT(secret string) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"exp": time.Now().Add(10 * time.Minute).Unix(), // InitState живёт 10 минут
//	})
//	return token.SignedString([]byte(secret))
//}
//
//func validateStateOauthJWT(stateToken, secret string) error {
//	token, err := jwt.Parse(stateToken, func(token *jwt.Token) (interface{}, error) {
//		return []byte(secret), nil
//	})
//	if err != nil {
//		return err
//	}
//	if !token.Valid {
//		return errors.New("token is not valid")
//	}
//
//	return nil
//}

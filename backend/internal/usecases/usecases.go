package usecases

import "github.com/saime-0/nice-pea-chat/internal/model"

type AuthUsecase interface {
	Auth(AuthIn) (AuthOut, error)
}

type AuthIn struct {
	Login string
}

type AuthOut struct {
	AccessToken string
}

type HealthcheckUsecase interface {
	Healthcheck() (HealthcheckOut, error)
}

type HealthcheckOut struct {
	MinVersionSupport string
	MinCodeSupport    int
	MaxVersionSupport string
	MaxCodeSupport    int
}

type UserByTokenUsecase interface {
	UserByToken(UserByTokenIn) (UserByTokenOut, error)
}

type UserByTokenIn struct {
	Token string
}

type UserByTokenOut struct {
	Found bool
	User  model.User
	Creds model.Credentials
}

type UserByIDUsecase interface {
	UserByID(UserByIDIn) (UserByIDOut, error)
}

type UserByIDIn struct {
	ID model.ID
}

type UserByIDOut struct {
	Found bool
	User  model.User
}

type UserUpdateUsecase interface {
	UserUpdate(UserUpdateIn) (UserUpdateOut, error)
}

type UserUpdateIn struct {
	Username string
}

type UserUpdateOut struct{}


type UserChatsUsecase struct {
	UserChats(UserChatsIn) (UserChatsOut, error)
}

type UserChatsIn struct {}
type UserChatsOut struct {
	
}

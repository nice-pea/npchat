package usecases

type AuthUsecase interface {
	Auth(AuthIn) AuthOut, error
}

type AuthIn struct {
	Login string
}

type AuthOut struct {
	AccessToken string
}

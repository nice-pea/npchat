package http

import (
	"net/http"

	ucAuthn "github.com/saime-0/nice-pea-chat/internal/usecase/authn"
	ucLogin "github.com/saime-0/nice-pea-chat/internal/usecase/login"
	ucPermissions "github.com/saime-0/nice-pea-chat/internal/usecase/permissions"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
	"github.com/saime-0/nice-pea-chat/internal/usecase/users"
	"github.com/saime-0/nice-pea-chat/internal/usecase/users/ucUserUpdate"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	m.handle("/health", Health)
	m.handle("/roles", Roles)
	m.handle("/users", Users)
	m.handle("/users/update", UserUpdate)
	m.handle("/permissions", Permissions)
	m.handle("/authn", Authn)
	m.handle("/login", Login)
}

func UserUpdate(req Request) (any, error) {
	ucParams := ucUserUpdate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.User); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Users(req Request) (any, error) {
	ucParams := users.Params{
		DB: req.DB,
	}

	return ucParams.Run()
}

func Login(req Request) (any, error) {
	ucParams := ucLogin.Params{
		Key: req.Form.Get("key"),
		DB:  req.DB,
	}

	return ucParams.Run()
}

func Authn(req Request) (any, error) {
	ucParams := ucAuthn.Params{
		Token: req.Form.Get("token"),
		DB:    req.DB,
	}

	return ucParams.Run()
}

func Permissions(req Request) (any, error) {
	ucParams := ucPermissions.Params{
		Locale: req.Locale,
		L10n:   req.L10n,
	}

	return ucParams.Run()
}

func Roles(req Request) (_ any, err error) {
	// Load params
	ucParams := ucRoles.Params{
		IDs:  nil,
		Name: req.Form.Get("name"),
		DB:   req.DB,
	}
	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Health(req Request) (any, error) {
	return req.L10n.Localize("none:ok", req.Locale, nil)
}

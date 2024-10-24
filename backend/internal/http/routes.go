package http

import (
	"net/http"

	ucAuth "github.com/saime-0/nice-pea-chat/internal/usecase/auth"
	ucPermissions "github.com/saime-0/nice-pea-chat/internal/usecase/permissions"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	m.handle("/health", Health)
	m.handle("/roles", Roles)
	m.handle("/permissions", Permissions)
	m.handle("/auth", Auth)
}

func Auth(req Request) (any, error) {
	ucParams := ucAuth.Params{
		Key: req.Form.Get("key"),
		DB:  req.DB,
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

package http

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/model"
	ucPermissions "github.com/saime-0/nice-pea-chat/internal/usecase/permissions"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	m.handle("/health", Health)
	m.handle("/roles", Roles)
	m.handle("/permissions", Permissions)
}

func Permissions(req Request) (_ any, err error) {
	ucParams := ucPermissions.Params{
		Locale: req.Locale,
		L10n:   req.L10n,
	}

	var perms []model.Permission
	if perms, err = ucParams.Run(); err != nil {
		return nil, err
	}

	return perms, nil
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

	var roles []model.Role
	if roles, err = ucParams.Run(); err != nil {
		return nil, err
	}

	return roles, nil
}

func Health(req Request) (any, error) {
	return req.L10n.Localize("none:ok", req.Locale, nil)
}

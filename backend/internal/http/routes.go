package http

import (
	"net/http"

	. "github.com/saime-0/nice-pea-chat/internal/model/role"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	m.handle("/healthz", Healthz)
	m.handle("/roles", Roles)
}

func Roles(req Request) (_ any, err error) {
	var ucParams ucRoles.Params

	ucParams.IDs, err = uintsParam(req.Form, "ids")
	if err != nil {
		return nil, err
	}

	ucParams.Name = req.Form.Get("name")

	var roles []Role
	if roles, err = ucParams.Run(); err != nil {
		return nil, err
	}

	return roles, nil
}

func Healthz(req Request) (any, error) {
	return req.L10n.Localize("none:ok", req.Locale, nil)
}

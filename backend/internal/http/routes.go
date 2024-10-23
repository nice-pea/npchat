package http

import (
	"net/http"
)

func declareRoutes(muxHttp *http.ServeMux) {
	m := mux{muxHttp}
	m.handle("/healthz", Healthz)
}

func Healthz(req Request) (any, error) {
	return req.L10n.Localize("none:ok", req.Locale, nil)
}

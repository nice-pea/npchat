package http

import (
	"net/http"
)

func declareRoutes(muxHttp *http.ServeMux) {
	m := mux{muxHttp}
	m.handle("/healthz", Healthz)
}

func Healthz(r *http.Request) (any, error) {
	return "ok", nil
}

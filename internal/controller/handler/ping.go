package handler

import "github.com/saime-0/nice-pea-chat/internal/controller/http2"

func Ping(router http2.Router) {
	router.HandleFunc("/ping", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

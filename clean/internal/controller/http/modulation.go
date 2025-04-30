package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func initContext(c *Controller, r *http.Request) Context {
	return Context{
		requestID: r.Header.Get("X-Request-ID"),
		subjectID: r.Header.Get("X-Subject-ID"), // TODO: переместить в auth
		request:   r,
	}
}

func (c *Controller) modulation(handle HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			respData any
			b        []byte
			err      error
		)

		// Получить значения из URL
		if err = r.ParseForm(); err != nil {
			slog.Warn("modulation: parse form: "+err.Error(),
				slog.String("url", r.RequestURI),
				slog.String("host", r.Host),
				slog.String("referer", r.Referer()),
			)
		}

		// Инициализация контекста запроса
		ctx := initContext(c, r)

		// Выполнить обработку запроса
		respData, err = handle(ctx)
		if err != nil {
			// Если есть ошибка
			w.WriteHeader(http.StatusBadRequest)
			respData = ResponseError{Error: err.Error()}
		}

		// Если ответ это строка, перезаписать структурой
		if s, ok := respData.(string); ok {
			// Если ответ это строка
			respData = ResponseMsg{Message: s}
		}

		// Сериализация ответа
		if b, err = json.Marshal(respData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b = []byte(fmt.Sprintf("{\"error\":\"%v\"}", err))
			slog.Error("modulation: marshal json: "+err.Error(),
				slog.String("url", r.RequestURI),
				slog.String("host", r.Host),
				slog.String("referer", r.Referer()),
			)
		}

		// Отправить ответ
		w.Header().Add("Content-Type", "application/json")
		if _, err = w.Write(b); err != nil {
			slog.Error("modulation: write bytes: "+err.Error(),
				slog.String("url", r.RequestURI),
				slog.String("host", r.Host),
				slog.String("referer", r.Referer()),
			)
		}
	}
}

type ResponseError struct {
	Error   string `json:"error"`
	ErrCode string `json:"errcode"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}

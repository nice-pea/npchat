package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var (
	ErrJsonMarshalResponseData = errors.New("json marshal response data")
	ErrWriteResponseBytes      = errors.New("write response bytes")
	ErrParseRequestURL         = errors.New("parse request url")
)

func (c *Controller) modulation(handle HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			respData any
			b        []byte
			err      error
		)

		// Получить значения из URL
		if err = r.ParseForm(); err != nil {
			logWarn(r, errors.Join(ErrParseRequestURL, err))
		}

		// Выполнить обработку запроса
		respData, err = handle(Context{request: r})
		if err != nil {
			// Если есть ошибка
			w.WriteHeader(httpStatusCodeByErr(err))
			respData = ResponseError{
				Error:   err.Error(),
				ErrCode: errCode(err),
			}
		}

		// Если ответ это строка, перезаписать структурой
		if s, ok := respData.(string); ok {
			// Если ответ это строка
			respData = ResponseMsg{Message: s}
		}

		// Сериализация ответа
		if b, err = json.Marshal(respData); err != nil {
			err = errors.Join(ErrJsonMarshalResponseData, err)
			w.WriteHeader(httpStatusCodeByErr(err))
			logErr(r, err)
			b = []byte(fmt.Sprintf(`{"error":"%v","errcode":"%v"}`, err, errCode(err)))
		}

		// Отправить ответ
		w.Header().Add("Content-Type", "application/json")
		if _, err = w.Write(b); err != nil {
			err = errors.Join(ErrWriteResponseBytes, err)
			w.WriteHeader(httpStatusCodeByErr(err))
			logErr(r, err)
		}
	}
}

func logErr(r *http.Request, err error) {
	slog.Error("modulation: "+err.Error(),
		slog.String("url", r.RequestURI),
		slog.String("host", r.Host),
		slog.String("referer", r.Referer()),
	)
}

func logWarn(r *http.Request, err error) {
	slog.Warn("modulation: "+err.Error(),
		slog.String("url", r.RequestURI),
		slog.String("host", r.Host),
		slog.String("referer", r.Referer()),
	)
}

type ResponseError struct {
	Error   string `json:"error"`
	ErrCode string `json:"errcode"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}

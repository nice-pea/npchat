package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Context представляет контекст HTTP-запроса
type Context struct {
	requestID string
	subjectID string
	//writer    http.ResponseWriter
	request *http.Request
}

type HandlerFunc func(Context) (any, error)

// Controller обрабатывает HTTP-запросы
type Controller struct {
	chats       service.Chats
	invitations service.Invitations
	members     service.Members

	mux *http.ServeMux
}

func InitController(chats service.Chats, invitations service.Invitations, members service.Members) *Controller {
	mux := http.NewServeMux()
	c := &Controller{
		chats:       chats,
		invitations: invitations,
		members:     members,
		mux:         mux,
	}

	mux.HandleFunc("POST /chats", c.modulation(c.CreateChat))
	mux.HandleFunc("GET /chats", c.modulation(c.GetChats))

	return c
}

func initContext(c *Controller, r *http.Request) Context {
	return Context{
		//Request: r.Request,
		//L10n:    s.L10n,
		//Locale:  locale(r.Header.Get("Accept-Language"), l10n.LocaleDefault),
		//Token:   getToken(r.Request),
		requestID: r.Header.Get("X-Request-ID"),
		subjectID: r.Header.Get("X-Subject-ID"), // TODO: переместить в auth
		//writer:    c,
		request: r,
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

// CreateChat создает новый чат
func (c *Controller) CreateChat(ctx Context) (any, error) {
	var input service.CreateInput
	if err := json.NewDecoder(ctx.request.Body).Decode(&input); err != nil {
		return nil, err
	}

	result, err := c.chats.Create(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetChats возвращает список чатов пользователя
func (c *Controller) GetChats(ctx Context) (any, error) {
	input := service.UserChatsInput{
		SubjectUserID: ctx.subjectID,
		UserID:        ctx.subjectID,
	}

	chats, err := c.chats.UserChats(input)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

// ServeHTTP обрабатывает HTTP-запросы
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

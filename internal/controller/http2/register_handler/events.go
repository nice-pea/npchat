package register_handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

// Events регистрирует обработчик для получения потока событий.
//
// Метод: GET /events
func Events(router *fiber.App, uc UsecasesForEvents, eventListener EventListener, jwtParser middleware.JwtParser) {
	router.Get(
		"/events",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			// Канал для обработки событий в отдельной горутине
			eventsChan := make(chan any)
			// Канал для обработки ошибок в отдельной горутине
			errorsChan := make(chan error)

			removeListener, err := eventListener.AddListener(UserID(ctx), SessionID(ctx), func(event events.Event, err error) {
				if err != nil {
					errorsChan <- err
				}
				if event.Type != "" {
					eventsChan <- event
				}
			})
			if err != nil {
				close(eventsChan)
				close(errorsChan)
				return err
			}

			// Канал для отслеживания завершения запроса
			reqCtxDone := ctx.Context().Done()

			// Регистрация обработчика для отправки потока сообщений
			ctx.Set("Content-Type", "text/event-stream")
			ctx.Set("Cache-Control", "no-cache")
			ctx.Set("Connection", "keep-alive")
			ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
				// Таймер для отправки keepalive
				keepAliveTickler := time.NewTicker(time.Second * 5)

				// Действия при звершении прослушивания событий
				defer func() {
					removeListener()
					keepAliveTickler.Stop()
					close(errorsChan)
					close(eventsChan)
				}()

				for {
					var errFprint error
					select {
					case event := <-eventsChan:
						// Отправлять события как json
						b, _ := json.Marshal(event)
						_, errFprint = fmt.Fprint(w, formatSSEMessage("event", string(b)))
					case err := <-errorsChan:
						// Отправлять ошибки
						_, errFprint = fmt.Fprint(w, formatSSEMessage("error", err.Error()))
					case <-keepAliveTickler.C:
						// Отправлять keepalive
						_, errFprint = fmt.Fprint(w, formatSSEMessage("keepalive", ""))
					case <-reqCtxDone:
						return
					}
					if errFprint != nil {
						slog.Warn("Events: SetBodyStreamWriter: fmt.Fprint: " + errFprint.Error())
						return
					}

					// Отправить данные во writer и очистить буфер
					if err := w.Flush(); err != nil {
						if err.Error() != "connection closed" {
							slog.Warn("Events: SetBodyStreamWriter: w.Flush: " + err.Error())
						}
						return
					}
				}
			})

			return nil
		},
	)
}

func formatSSEMessage(eventType, data string) string {
	sb := strings.Builder{}

	if eventType != "" {
		sb.WriteString(fmt.Sprintf("event: %s\n", eventType))
	}
	sb.WriteString(fmt.Sprintf("retry: %d\n", 15000))
	if data != "" {
		sb.WriteString(fmt.Sprintf("data: %v\n", data))
	}
	sb.WriteString("\n")

	return sb.String()
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForEvents interface {
	middleware.UsecasesForRequireAuthorizedSession
}

type EventListener interface {
	AddListener(userID, sessionID uuid.UUID, f func(event events.Event, err error)) (removeListener func(), err error)
}

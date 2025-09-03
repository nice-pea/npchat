package register_handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

// Events регистрирует обработчик для получения потока событий.
//
// Метод: GET /events
func Events(router *fiber.App, uc UsecasesForEvents, eventListener EventListener) {
	router.Get(
		"/events",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
		func(ctx *fiber.Ctx) error {
			ctx.Set("Content-Type", "text/event-stream")
			ctx.Set("Cache-Control", "no-cache")
			ctx.Set("Connection", "keep-alive")
			// ctx.Set("Transfer-Encoding", "chunked")

			session := Session(ctx)

			// Таймер для отправки keepalive
			keepAliveTickler := time.NewTicker(time.Second * 5)
			// Контекст для корректной остановки прослушивания событий
			ctx2, cancel := context.WithCancel(context.Background())
			// Канал для отслеживания завершения запроса
			reqCtxDone := ctx.Context().Done()

			// Канал для обработки событий в отдельной горутине
			eventsChan := make(chan any)
			// Канал для обработки ошибок в отдельной горутине
			errorsChan := make(chan error)

			// Регистрация обработчика для отправки потока сообщений
			ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				go func() {
					defer func() {
						cancel()
						keepAliveTickler.Stop()
						close(eventsChan)
					}()

					for {
						select {
						case event := <-eventsChan:
							b, _ := json.Marshal(event)
							fmt.Fprint(w, formatSSEMessage("event", string(b)))
						case err := <-errorsChan:
							fmt.Fprint(w, formatSSEMessage("error", err.Error()))
						case <-keepAliveTickler.C:
							fmt.Fprint(w, formatSSEMessage("keepalive", ""))
						case <-reqCtxDone:
							return
						}
						// Отправить данные во writer и очистить буфер
						if err := w.Flush(); err != nil {
							slog.Warn("Events: w.Flush: " + err.Error())
							return
						}
					}
				}()
				err := eventListener.Listen(ctx2, session.UserID, session.ID, func(event any) {
					eventsChan <- event
				})
				if err != nil && !errors.Is(err, context.Canceled) {
					errorsChan <- err
				}
			}))

			return nil
		},
	)
}

func formatSSEMessage(eventType, data string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", eventType))
	sb.WriteString(fmt.Sprintf("retry: %d\n", 15000))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", data))

	return sb.String()
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForEvents interface {
	middleware.UsecasesForRequireAuthorizedSession
}

type EventListener interface {
	Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error
}

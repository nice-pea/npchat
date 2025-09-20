package register_handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	myChats "github.com/nice-pea/npchat/internal/usecases/chats/my_chats"
)

// MyChats регистрирует HTTP-обработчик для получения списка чатов пользователя.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func MyChats(router *fiber.App, uc UsecasesForMyChats, jparser middleware.JwtParser) {
	router.Get(
		"/chats",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
		func(ctx *fiber.Ctx) error {
			pageToken, err := decodeKeyset(ctx)
			if err != nil {
				return err
			}

			input := myChats.In{
				SubjectID: UserID(ctx),
				UserID:    UserID(ctx),
				PageToken: pageToken,
			}

			out, err := uc.MyChats(input)
			if err != nil {
				return err
			}
			nextPageToken, err := encodeKeyset(ctx, out.NextKeyset)
			if err != nil {
				return err
			}

			return ctx.JSON(fiber.Map{
				"Chats":           out.Chats,
				"next_page_token": nextPageToken,
			})
		},
	)
}

// UsecasesForMyChats определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForMyChats interface {
	MyChats(myChats.In) (myChats.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}

func decodeKeyset(ctx *fiber.Ctx) (myChats.Keyset, error) {
	b, err := base64.StdEncoding.DecodeString(ctx.Query("page_token"))
	if err != nil {
		return myChats.Keyset{}, fmt.Errorf("decode page_token: %w", err)
	}
	var keyset myChats.Keyset
	return keyset, json.Unmarshal(b, &keyset)
}

func encodeKeyset(ctx *fiber.Ctx, keyset myChats.Keyset) (string, error) {
	b, err := json.Marshal(keyset)
	if err != nil {
		return "", fmt.Errorf("marshal keyset: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

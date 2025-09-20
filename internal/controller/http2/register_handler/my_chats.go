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
func MyChats(router *fiber.App, uc UsecasesForMyChats, jwtParser middleware.JwtParser) {
	router.Get(
		"/chats",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			keyset, err := decodeKeyset(ctx.Query("page_token"))
			if err != nil {
				return err
			}

			input := myChats.In{
				SubjectID: UserID(ctx),
				UserID:    UserID(ctx),
				Keyset:    keyset,
			}

			out, err := uc.MyChats(input)
			if err != nil {
				return err
			}
			nextPageToken, err := encodeKeyset(out.NextKeyset)
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

// decodeKeyset расшифровывает строку в формате base64 и разбирает ее на Keyset
func decodeKeyset(pageToken string) (myChats.Keyset, error) {
	if pageToken == "" {
		return myChats.Keyset{}, nil
	}
	b, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return myChats.Keyset{}, fmt.Errorf("decode page_token: %w", err)
	}
	var keyset myChats.Keyset
	return keyset, json.Unmarshal(b, &keyset)
}

// encodeKeyset преобразует keyset в json и кодирует в строку в base64
func encodeKeyset(keyset myChats.Keyset) (string, error) {
	if keyset == (myChats.Keyset{}) {
		return "", nil
	}
	b, err := json.Marshal(keyset)
	if err != nil {
		return "", fmt.Errorf("marshal keyset: %w", err)
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

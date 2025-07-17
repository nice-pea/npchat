package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

func Session(context *fiber.Ctx) sessionn.Session {
	return context.Locals(middleware.CtxUserSession, sessionn.Session{}).(sessionn.Session)
}

// ParamsUUID возвращает значение из пути запроса как uuid
func ParamsUUID(context *fiber.Ctx, name string) uuid.UUID {
	val, _ := uuid.Parse(context.Params(name))
	return val
}

package register_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

func UserID(context *fiber.Ctx) uuid.UUID {
	return context.Locals(middleware.CtxKeyUserID).(uuid.UUID)
}
func SessionID(context *fiber.Ctx) uuid.UUID {
	return context.Locals(middleware.CtxKeySessionID).(uuid.UUID)
}

// ParamsUUID возвращает значение из пути запроса как uuid
func ParamsUUID(context *fiber.Ctx, name string) uuid.UUID {
	val, _ := uuid.Parse(context.Params(name))
	return val
}

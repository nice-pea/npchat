package register_handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

// loginResultData возвращает результат регистрации входа в виде Map
func loginResultData(session sessionn.Session, user userr.User, issuer JwtIssuer) fiber.Map {
	data := fiber.Map{
		"session": session,
		"user":    user,
	}

	// Если есть JWT-issuer, то генерируем JWT
	if issuer != nil {
		var err error
		if data["jwt"], err = issuer.Issue(session); err != nil {
			slog.Error("Ошибка генерации JWT", "error", err)
		}
	}

	return data
}

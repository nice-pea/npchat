package service 

import (
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

// Users сервис, объединяющий случаи использования(юзкейсы) в контексте агрегата пользователей
type Users struct {
	Providers    OAuthProviders // Карта провайдеров OAuth
	Repo         userr.Repository
	SessionsRepo sessionn.Repository // Репозиторий сессий пользователей
}

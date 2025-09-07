package register_handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockRegisterHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler/mocks"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	oauthAuthorize "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_authorize"
	oauthComplete "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_complete"
)

func TestOauthAuthorize(t *testing.T) {
	t.Run("редирект после успешного выполнения", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		// Настройка мока
		const redirectURL = "https://provider.example/oauth"
		mockUsecases := mockRegisterHandler.NewUsecasesForOauthAuthorize(t)
		mockUsecases.
			On("OauthAuthorize", mock.Anything).
			Return(oauthAuthorize.Out{
				RedirectURL: redirectURL,
				State:       "test-state",
			}, nil)

		// Регистрация обработчика
		OauthAuthorize(fiberApp, mockUsecases)

		// Выполнить запрос
		req := httptest.NewRequest("GET", "/oauth/someProvider/authorize", nil)
		resp, err := fiberApp.Test(req)
		assert.NoError(t, err)
		// Удостовериться, что сервер сделает редирект
		assert.Equal(t, fiber.StatusTemporaryRedirect, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Location"), redirectURL)

		// Проверить куку oauth state
		cookies := resp.Cookies()
		assert.NotEmpty(t, cookies, "должна быть установлена кука")
		cookie := cookies[0]
		assert.Equal(t, "someProvider-oauthState", cookie.Name)
		assert.Equal(t, "test-state", cookie.Value)
		assert.True(t, cookie.Secure, "кука должна быть secure")
		assert.True(t, cookie.HttpOnly, "кука должна быть httpOnly")

	})

	t.Run("не ОК при ошибке в юзкейсе", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		// Настройка мока
		mockUsecases := mockRegisterHandler.NewUsecasesForOauthAuthorize(t)
		mockUsecases.
			On("OauthAuthorize", mock.Anything).
			Return(oauthComplete.Out{}, errors.New("some error"))

		// Регистрация обработчика
		OauthAuthorize(fiberApp, mockUsecases)

		// Выполнить запрос
		req := httptest.NewRequest("GET", "/oauth/someProvider/authorize", nil)
		resp, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)

		// Кук нет
		assert.Empty(t, resp.Cookies())
	})
}

func TestOauthCallback(t *testing.T) {
	t.Run("успешное выполнение", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		// Настройка мока
		mockUsecases := mockRegisterHandler.NewUsecasesForOauthCallback(t)
		mockUsecases.
			On("OauthComplete", mock.Anything).
			Return(oauthComplete.Out{
				Session: sessionn.Session{ID: uuid.New()},
				User:    userr.User{ID: uuid.New()},
			}, nil)

		// Регистрация обработчика
		OauthCallback(fiberApp, mockUsecases)

		// Создать запрос
		req := httptest.NewRequest("GET", "/oauth/someProvider/callback?state=someValue", nil)
		// Установить куки
		req.AddCookie(&http.Cookie{
			Name:  "someProvider-oauthState",
			Value: "someValue",
		})
		// Выполнить запрос
		resp, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Должен содержаться ответ
		b, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.NotEmpty(t, b)
	})

	t.Run("не ОК при отсутствии кук", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		// Настройка мока
		mockUsecases := mockRegisterHandler.NewUsecasesForOauthCallback(t)

		// Регистрация обработчика
		OauthCallback(fiberApp, mockUsecases)

		// Создать запрос
		req := httptest.NewRequest("GET", "/oauth/someProvider/callback?state=someValue", nil)
		// Выполнить запрос
		resp, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("не ОК при ошибке в юзкейсе", func(t *testing.T) {
		fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
		// Настройка мока
		mockUsecases := mockRegisterHandler.NewUsecasesForOauthCallback(t)
		mockUsecases.
			On("OauthComplete", mock.Anything).
			Return(oauthComplete.Out{}, errors.New("some error"))

		// Регистрация обработчика
		OauthCallback(fiberApp, mockUsecases)

		// Создать запрос
		req := httptest.NewRequest("GET", "/oauth/someProvider/callback?state=someValue", nil)
		// Установить куки
		req.AddCookie(&http.Cookie{
			Name:  "someProvider-oauthState",
			Value: "someValue",
		})
		// Выполнить запрос
		resp, err := fiberApp.Test(req)
		assert.NoError(t, err)
		assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)
	})
}

package http2

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type responseMessage struct {
	Message string `json:"message"`
}

type responseError struct {
	Error   string `json:"error"`
	ErrCode string `json:"errcode"`
}

// TestPing tests the Ping handler
func (suite *controllerTestSuite) TestPing() {
	suite.Run("на ping вернется pong", func() {
		resp, err := http.Get(suite.server.URL + "/ping")
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Require().Equal(http.StatusOK, resp.StatusCode)

		var respData responseMessage
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)

		suite.Equal(responseMessage{Message: "pong"}, respData)
	})
}

const headerXRequestID = "X-Request-ID"
const headerAccept = "Accept"
const headerAuthorization = "Authorization"

func (suite *controllerTestSuite) newClientRequest(method, path, token string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, suite.server.URL+path, body)
	suite.Require().NoError(err)
	req.Header.Set(headerXRequestID, uuid.NewString())
	req.Header.Set(headerAccept, "application/json")
	if token != "" {
		req.Header.Set(headerAuthorization, "Bearer "+token)
	}
	return req
}

func (suite *controllerTestSuite) newRndUserWithSession(sessionStatus int) (out struct {
	User    domain.User
	Session domain.Session
}) {
	out.User = domain.User{
		ID: uuid.NewString(),
	}
	err := suite.rr.users.Save(out.User)
	suite.Require().NoError(err)

	out.Session = domain.Session{
		ID:     uuid.NewString(),
		UserID: out.User.ID,
		Token:  uuid.NewString(),
		Status: sessionStatus,
	}
	err = suite.rr.sessions.Save(out.Session)
	suite.Require().NoError(err)

	return
}

func (suite *controllerTestSuite) newRndUserWithAuthnPassword() domain.AuthnPassword {
	ap := domain.AuthnPassword{
		UserID:   uuid.NewString(),
		Login:    uuid.NewString(),
		Password: uuid.NewString(),
	}
	err := suite.rr.users.Save(domain.User{ID: ap.UserID})
	suite.Require().NoError(err)
	err = suite.rr.authnPassword.Save(ap)
	suite.Require().NoError(err)

	return ap
}
func (suite *controllerTestSuite) jsonBody(v any) io.Reader {
	body, err := json.Marshal(v)
	suite.Require().NoError(err)
	return bytes.NewBuffer(body)
}

// saveChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *controllerTestSuite) saveChat(chat domain.Chat) domain.Chat {
	err := suite.rr.chats.Save(chat)
	suite.Require().NoError(err)

	return chat
}

// saveMember сохраняет участника в репозиторий, в случае ошибки завершит тест
func (suite *controllerTestSuite) saveMember(member domain.Member) domain.Member {
	err := suite.rr.members.Save(member)
	suite.Require().NoError(err)

	return member
}

func (suite *controllerTestSuite) TestClientMiddlewares() {
	const existingClientAPIEndpoint = "/chats"

	suite.Run("идентификатор запроса/обязательно должен быть передан заголовок с идентификатором запроса", func() {
		// Создать запрос
		req, err := http.NewRequest("GET", suite.server.URL+existingClientAPIEndpoint, nil)
		suite.Require().NoError(err)

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		// Проверить код ответа
		suite.Require().Equal(http.StatusBadRequest, resp.StatusCode)

		// Проверить ответ
		var respData responseError
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)
		suite.Equal(ErrUnknownRequestID.Error(), respData.Error)
		suite.Equal(ErrCodeInvalidXRequestIDHeader, respData.ErrCode)
	})

	suite.Run("тип содержимого/api поддерживает только контент json", func() {
		// Создать запрос
		req, err := http.NewRequest("GET", suite.server.URL+existingClientAPIEndpoint, nil)
		suite.Require().NoError(err)
		req.Header.Set(headerXRequestID, uuid.NewString())

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		// Проверить код ответа
		suite.Require().Equal(http.StatusBadRequest, resp.StatusCode)

		// Проверить ответ
		var respData responseError
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)
		suite.Equal(ErrUnsupportedAcceptedContentType.Error(), respData.Error)
		suite.Equal(ErrCodeUnsupportedAcceptedContentType, respData.ErrCode)
	})

	suite.Run("аутентификация сессии/защищенные эндпоинты будут возвращать ошибку", func() {
		// Создать запрос
		req, err := http.NewRequest("GET", suite.server.URL+existingClientAPIEndpoint, nil)
		suite.Require().NoError(err)
		req.Header.Set(headerXRequestID, uuid.NewString())
		req.Header.Set(headerAccept, "application/json")

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		// Проверить код ответа
		suite.Require().Equal(http.StatusUnauthorized, resp.StatusCode)

		// Проверить ответ
		var respData responseError
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)
		suite.Equal(ErrUnauthorized.Error(), respData.Error)
		suite.Equal(ErrCodeInvalidAuthorizationHeader, respData.ErrCode)
	})

	suite.Run("аутентификация сессии/запросы Verified сессий проходят успешно", func() {
		// Создать запрос
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		req := suite.newClientRequest("GET", existingClientAPIEndpoint, uws.Session.Token, nil)

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		// Проверить код ответа
		if !suite.Equal(http.StatusOK, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})
}

func (suite *controllerTestSuite) TestLoginByPassword() {
	suite.Run("успешная авторизация", func() {
		// Создать нового пользователя с AuthnPassword
		ap := suite.newRndUserWithAuthnPassword()

		// Создать запрос
		body := suite.jsonBody(map[string]string{
			"login":    ap.Login,
			"password": ap.Password,
		})
		req := suite.newClientRequest("POST", "/login/password", "", body)

		// Выполнение запроса
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		// Проверка результата
		if !suite.Equal(http.StatusOK, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})

	suite.Run("неверные учетные данные", func() {
		// Создать нового пользователя с AuthnPassword
		lc := suite.newRndUserWithAuthnPassword()

		credentials := map[string]string{
			"login":    lc.Login,
			"password": "wrongpass",
		}
		body, err := json.Marshal(credentials)
		suite.Require().NoError(err)

		req := suite.newClientRequest("POST", "/login/password", "", bytes.NewBuffer(body))

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		if !suite.Equal(http.StatusBadRequest, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})
}

func (suite *controllerTestSuite) TestMyChats() {
	suite.Run("получение списка чатов", func() {
		// Создаем тестового пользователя с сессией
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)

		req := suite.newClientRequest("GET", "/chats", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		if !suite.Equal(http.StatusOK, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})
}

func (suite *controllerTestSuite) TestCreateChat() {
	suite.Run("создание нового чата", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)

		// Создать запрос
		body := suite.jsonBody(map[string]string{
			"name":          "somechat",
			"chief_user_id": uws.User.ID,
		})
		req := suite.newClientRequest("POST", "/chats", uws.Session.Token, body)

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		if !suite.Equal(http.StatusOK, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})
}

func (suite *controllerTestSuite) TestChatMembers() {
	suite.Run("получение списка участников чата", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)

		// Создать чат
		chat := suite.saveChat(domain.Chat{ID: uuid.NewString()})
		// Создать участников
		for range 10 {
			suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
		}
		// участник Subject
		suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uws.User.ID,
			ChatID: chat.ID,
		})

		// Создать запрос
		req := suite.newClientRequest("GET", "/chats/"+chat.ID+"/members", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		if !suite.Equal(http.StatusOK, resp.StatusCode) {
			suite.PrintBody(resp)
		}
	})
}

func (suite *controllerTestSuite) PrintBody(resp *http.Response) {
	b, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	log.Println(string(b))
}

func (suite *controllerTestSuite) TestUpdateChatName() {
	suite.Run("обновление названия чата", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		chatID := uuid.New().String()

		updateData := map[string]string{
			"name": "New Chat Name",
		}
		body, err := json.Marshal(updateData)
		suite.Require().NoError(err)

		req := suite.newClientRequest("PUT", "/chats/"+chatID+"/name", uws.Session.Token, bytes.NewBuffer(body))

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestChatInvitations() {
	suite.Run("получение списка приглашений в чат", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		chatID := uuid.New().String()

		req := suite.newClientRequest("GET", "/chats/"+chatID+"/invitations", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestMyInvitations() {
	suite.Run("получение списка моих приглашений", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)

		req := suite.newClientRequest("GET", "/invitations", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestSendInvitation() {
	suite.Run("отправка приглашения в чат", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)

		req := suite.newClientRequest("POST", "/invitations", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestAcceptInvitation() {
	suite.Run("принятие приглашения в чат", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		invitationID := uuid.New().String()

		req := suite.newClientRequest("POST", "/invitations/"+invitationID+"/accept", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestCancelInvitation() {
	suite.Run("отмена приглашения в чат", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		invitationID := uuid.New().String()

		req := suite.newClientRequest("POST", "/invitations/"+invitationID+"/cancel", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestLeaveChat() {
	suite.Run("выход из чата", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		chatID := uuid.New().String()

		req := suite.newClientRequest("POST", "/chats/"+chatID+"/leave", uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

func (suite *controllerTestSuite) TestDeleteMember() {
	suite.Run("удаление участника из чата", func() {
		uws := suite.newRndUserWithSession(domain.SessionStatusVerified)
		chatID := uuid.New().String()
		memberID := uuid.New().String()

		req := suite.newClientRequest("DELETE", "/chats/"+chatID+"/members/"+memberID, uws.Session.Token, nil)

		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer func() { _ = resp.Body.Close() }()

		suite.Equal(http.StatusOK, resp.StatusCode)
	})
}

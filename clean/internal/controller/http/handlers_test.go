package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
		defer resp.Body.Close() //nolint:errcheck

		suite.Require().Equal(http.StatusOK, resp.StatusCode)

		var respData responseMessage
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)

		suite.Equal(responseMessage{Message: "pong"}, respData)
	})
}

const headerXRequestID = "X-Request-ID"
const headerAccept = "Accept"

func (suite *controllerTestSuite) TestClientMiddlewares() {
	const existingClientAPIEndpoint = "/chats"

	suite.Run("идентификатор запроса/обязательно должен быть передан заголовок с идентификатором запроса", func() {
		// Создать запрос
		req, err := http.NewRequest("GET", suite.server.URL+existingClientAPIEndpoint, nil)
		suite.Require().NoError(err)

		// Выполнить запрос
		resp, err := http.DefaultClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close() //nolint:errcheck

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
		defer resp.Body.Close() //nolint:errcheck

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
		defer resp.Body.Close() //nolint:errcheck

		// Проверить код ответа
		suite.Require().Equal(http.StatusUnauthorized, resp.StatusCode)

		// Проверить ответ
		var respData responseError
		err = json.NewDecoder(resp.Body).Decode(&respData)
		suite.Require().NoError(err)
		suite.Equal(ErrUnauthorized.Error(), respData.Error)
		suite.Equal(ErrCodeInvalidAuthorizationHeader, respData.ErrCode)
	})
}

// TestCreateChat tests the CreateChat handler
func (suite *controllerTestSuite) TestCreateChat() {
	suite.Run("should create a chat successfully", func() {
		// This handler should:
		// 1. Parse the request body to get the chat name and chief user ID
		// 2. Call the chats.Create service method
		// 3. Return the created chat and chief member
		suite.T().Skip("Implementation exists, test needs to be implemented")
	})
}

// TestGetChats tests the GetChats handler
func (suite *controllerTestSuite) TestGetChats() {
	suite.Run("should return user chats", func() {
		// This handler should:
		// 1. Get the user ID from the session
		// 2. Call the chats.UserChats service method
		// 3. Return the list of chats
		suite.T().Skip("Implementation exists, test needs to be implemented")
	})
}

// Tests for not yet implemented handlers - Chats service

// TestUpdateChatName tests the UpdateChatName handler
func (suite *controllerTestSuite) TestUpdateChatName() {
	suite.Run("should update chat name", func() {
		// This handler should:
		// 1. Parse the request body to get the chat ID and new name
		// 2. Get the user ID from the session
		// 3. Call the chats.UpdateName service method
		// 4. Return the updated chat
		suite.T().Skip("Handler not yet implemented")
	})
}

// Tests for not yet implemented handlers - Invitations service

// TestChatInvitations tests the ChatInvitations handler
func (suite *controllerTestSuite) TestChatInvitations() {
	suite.Run("should return chat invitations", func() {
		// This handler should:
		// 1. Get the chat ID from the URL path
		// 2. Get the user ID from the session
		// 3. Call the invitations.ChatInvitations service method
		// 4. Return the list of invitations
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestUserInvitations tests the UserInvitations handler
func (suite *controllerTestSuite) TestUserInvitations() {
	suite.Run("should return user invitations", func() {
		// This handler should:
		// 1. Get the user ID from the session
		// 2. Call the invitations.UserInvitations service method
		// 3. Return the list of invitations
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestSendInvitation tests the SendInvitation handler
func (suite *controllerTestSuite) TestSendInvitation() {
	suite.Run("should send invitation", func() {
		// This handler should:
		// 1. Parse the request body to get the chat ID and user ID
		// 2. Get the subject user ID from the session
		// 3. Call the invitations.SendInvitation service method
		// 4. Return the created invitation
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestAcceptInvitation tests the AcceptInvitation handler
func (suite *controllerTestSuite) TestAcceptInvitation() {
	suite.Run("should accept invitation", func() {
		// This handler should:
		// 1. Parse the request body to get the chat ID
		// 2. Get the user ID from the session
		// 3. Call the invitations.AcceptInvitation service method
		// 4. Return success
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestCancelInvitation tests the CancelInvitation handler
func (suite *controllerTestSuite) TestCancelInvitation() {
	suite.Run("should cancel invitation", func() {
		// This handler should:
		// 1. Parse the request body to get the chat ID and user ID
		// 2. Get the subject user ID from the session
		// 3. Call the invitations.CancelInvitation service method
		// 4. Return success
		suite.T().Skip("Handler not yet implemented")
	})
}

// Tests for not yet implemented handlers - Members service

// TestChatMembers tests the ChatMembers handler
func (suite *controllerTestSuite) TestChatMembers() {
	suite.Run("should return chat members", func() {
		// This handler should:
		// 1. Get the chat ID from the URL path
		// 2. Get the user ID from the session
		// 3. Call the members.ChatMembers service method
		// 4. Return the list of members
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestLeaveChat tests the LeaveChat handler
func (suite *controllerTestSuite) TestLeaveChat() {
	suite.Run("should leave chat", func() {
		// This handler should:
		// 1. Get the chat ID from the URL path
		// 2. Get the user ID from the session
		// 3. Call the members.LeaveChat service method
		// 4. Return success
		suite.T().Skip("Handler not yet implemented")
	})
}

// TestDeleteMember tests the DeleteMember handler
func (suite *controllerTestSuite) TestDeleteMember() {
	suite.Run("should delete member", func() {
		// This handler should:
		// 1. Get the chat ID from the URL path
		// 2. Parse the request body to get the user ID to delete
		// 3. Get the subject user ID from the session
		// 4. Call the members.DeleteMember service method
		// 5. Return success
		suite.T().Skip("Handler not yet implemented")
	})
}

// Tests for not yet implemented handlers - Sessions service

// TestFindSession tests the FindSession handler
func (suite *controllerTestSuite) TestFindSession() {
	suite.Run("should find session", func() {
		// This handler should:
		// 1. Get the token from the query parameters
		// 2. Call the sessions.Find service method
		// 3. Return the found sessions
		suite.T().Skip("Handler not yet implemented")
	})
}

package http

import (
	"net/http"
)

func (suite *servicesTestSuite) Test_Server() {
	suite.Run("это http сервер", func() {
		var _ http.Handler = new(Controller)
	})
	suite.Run("context с доступом к данным", func() {
		c := &Context{
			requestID: "",
			subjectID: "",
			writer:    (http.ResponseWriter)(nil),
			request:   (*http.Request)(nil),
		}

		//var _ interface {
		//	RequestID() string
		//	SubjectID() string
		//	Writer() http.ResponseWriter
		//	Request() http.Request
		//} = c
		_ = c
	})
	suite.Run("есть метод modulation", func() {
		var _ func(http.Handler) func(Context) (any, error) = modulation
		//modulation(http.Handler) func() (any, error)
	})
	suite.Run("наличие методов", func() {
		var _ interface {
			CreateChat(Context) (any, error)
			GetChats(Context) (any, error)
		} = new(Controller)
	})

	//c := Controller{
	//	chats:       service.Chats{},
	//	invitations: service.Invitations{},
	//	members:     service.Members{},
	//}
	//http.NewServeMux()
	//
	//require.IsType(t, new(interface {
	//	chats() service.Chats
	//	invitations() service.Invitations
	//	members() service.Members
	//}), c)
	//
	//addr := ":8080"
	//go func() {
	//	err := http.ListenAndServe(addr, c)
	//	assert.NoError(t, err)
	//}()

}

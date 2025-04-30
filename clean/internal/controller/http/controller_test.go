package http

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *controllerTestSuite) Test_Server() {
	suite.Run("это http сервер", func() {
		var _ http.Handler = new(Controller)
	})
	suite.Run("context с доступом к данным", func() {
		_ = &Context{
			requestID: "",
			request:   (*http.Request)(nil),
			session:   domain.Session{},
		}
	})
	suite.Run("есть метод modulation", func() {
		//// Функция для преобразования controller.HandlerFunc в тип http.HandlerFunc
		//var f func(HandlerFunc) http.HandlerFunc
		//// Проверить существование такой функции
		//f = modulation
		//_ = f
		// Функция для преобразования controller.HandlerFunc в тип http.HandlerFunc
		var _ interface {
			modulation(HandlerFunc) http.HandlerFunc
		} = new(Controller)
	})
	//suite.Run("есть middleware для наполнения Context значениями", func() {
	//	// Функция для преобразования controller.HandlerFunc в тип http.HandlerFunc
	//	var mw func(HandlerFunc) HandlerFunc
	//	// Проверить существование такой функции
	//	mw = initContext
	//	_ = mw
	//})
	suite.Run("наличие методов", func() {
		// Контроллер содержит методы для обработки следующих запросов:
		var _ interface {
			// Создать чат
			CreateChat(Context) (any, error)
			// Получить список чатов
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

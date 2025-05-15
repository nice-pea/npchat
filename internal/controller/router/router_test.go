package router

//
//import (
//	"net/http"
//
//	"github.com/saime-0/nice-pea-chat/internal/domain"
//)
//
//func (suite *controllerTestSuite) Test_Server() {
//	suite.Run("это http2 сервер", func() {
//		var _ http.Handler = new(Router)
//	})
//	suite.Run("context с доступом к данным", func() {
//		_ = &Context{
//			requestID: "",
//			request:   (*http.Request)(nil),
//			session:   domain.Session{},
//		}
//	})
//	suite.Run("есть метод modulation", func() {
//		// Функция для преобразования controller.HandlerFunc в тип http2.HandlerFunc
//		var _ interface {
//			modulation(HandlerFunc) http.HandlerFunc
//		} = new(Router)
//	})
//	suite.Run("наличие методов", func() {
//		// Контроллер содержит методы для обработки следующих запросов:
//		var _ interface {
//			LoginByPassword(Context) (any, error)  // Авторизация по логину/паролю
//			MyChats(Context) (any, error)          // Получить список чатов
//			CreateChat(Context) (any, error)       // Создать чат
//			ChatMembers(Context) (any, error)      // Получить список участников чата
//			UpdateChatName(Context) (any, error)   // Обновить название чата
//			ChatInvitations(Context) (any, error)  // Получить список приглашений в чат
//			MyInvitations(Context) (any, error)    // Получить список моих приглашений
//			SendInvitation(Context) (any, error)   // Отправить приглашение в чат
//			AcceptInvitation(Context) (any, error) // Принять приглашение в чат
//			CancelInvitation(Context) (any, error) // Отменить приглашение в чат
//			LeaveChat(Context) (any, error)        // Покинуть чат
//			DeleteMember(Context) (any, error)     // Удалить участника из чата
//		} = new(Router)
//	})
//}

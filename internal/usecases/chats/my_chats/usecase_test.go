package myChats

import (
	"testing"
	"time"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	mockChatt "github.com/nice-pea/npchat/internal/domain/chatt/mocks"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.SuiteWithMocks
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_Chats_MyChats() {
	suite.Run("запрещает просмотр чужих чатов", func() {
		usecase, _ := newUsecase(suite)

		input := In{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		}

		out, err := usecase.MyChats(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(out.Chats)
	})

	suite.Run("возвращает пустой список когда чатов нет", func() {
		usecase, mockRepo := newUsecase(suite)

		input := suite.newUserChatsInput(uuid.New())
		mockRepo.EXPECT().List(chatt.Filter{
			ParticipantID: input.UserID,
			ActiveBefore:  input.Keyset.ActiveBefore,
			Limit:         defaultPageSize,
		}).Return([]chatt.Chat{}, nil).Once()

		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Empty(out.Chats)
		suite.True(out.NextKeyset.ActiveBefore.IsZero())
	})

	suite.Run("возвращает чаты пользователя", func() {
		usecase, mockRepo := newUsecase(suite)

		userID := uuid.New()
		now := time.Now().UTC().Truncate(time.Millisecond)
		const countChats = 3
		expectedChats := make([]chatt.Chat, 0, countChats)
		for i := range countChats {
			chat := suite.RndChat()
			chat.LastActiveAt = now.Add(-time.Duration(i) * time.Minute)
			chat.Participants = append(chat.Participants, suite.NewParticipant(userID))
			expectedChats = append(expectedChats, chat)
		}

		input := suite.newUserChatsInput(userID)
		mockRepo.EXPECT().List(chatt.Filter{
			ParticipantID: userID,
			ActiveBefore:  input.Keyset.ActiveBefore,
			Limit:         defaultPageSize,
		}).Return(expectedChats, nil).Once()

		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Equal(expectedChats, out.Chats)
	})

	suite.Run("не возвращает keyset если элементов меньше лимита", func() {
		usecase, mockRepo := newUsecase(suite)

		userID := uuid.New()
		now := time.Now().UTC().Truncate(time.Millisecond)
		chats := make([]chatt.Chat, defaultPageSize-1)
		for i := range chats {
			chats[i] = suite.RndChat()
			chats[i].LastActiveAt = now.Add(-time.Duration(i) * time.Minute)
		}

		input := suite.newUserChatsInput(userID)
		mockRepo.EXPECT().List(chatt.Filter{
			ParticipantID: userID,
			ActiveBefore:  input.Keyset.ActiveBefore,
			Limit:         defaultPageSize,
		}).Return(chats, nil).Once()

		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.True(out.NextKeyset.ActiveBefore.IsZero())
	})

	suite.Run("учитывает фильтрацию по активности", func() {
		usecase, mockRepo := newUsecase(suite)

		userID := uuid.New()
		activeBefore := time.Now().Add(-time.Minute).UTC().Truncate(time.Millisecond)
		expectedChats := make([]chatt.Chat, 5)
		for i := range expectedChats {
			expectedChats[i] = suite.RndChat()
			expectedChats[i].LastActiveAt = activeBefore.Add(-time.Duration(i+1) * time.Minute)
		}

		input := In{
			SubjectID: userID,
			UserID:    userID,
			Keyset: Keyset{
				ActiveBefore: activeBefore,
			},
		}
		mockRepo.EXPECT().List(chatt.Filter{
			ParticipantID: userID,
			ActiveBefore:  activeBefore,
			Limit:         defaultPageSize,
		}).Return(expectedChats, nil).Once()

		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Equal(expectedChats, out.Chats)
		suite.True(out.NextKeyset.ActiveBefore.IsZero())
	})
}

func (suite *testSuite) newUserChatsInput(userID uuid.UUID) In {
	return In{
		SubjectID: userID,
		UserID:    userID,
	}
}

func newUsecase(suite *testSuite) (*MyChatsUsecase, *mockChatt.Repository) {
	uc := &MyChatsUsecase{
		Repo: suite.RR.Chats,
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	return uc, mockRepo
}

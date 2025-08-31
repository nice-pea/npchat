package deleteMember

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID          = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID             = errors.New("некорректное значение ChatID")
	ErrInvalidUserID             = errors.New("некорректное значение UserID")
	ErrMemberCannotDeleteHimself = errors.New("участник не может удалить самого себя")
	ErrSubjectUserIsNotChief     = errors.New("пользователь не является главным администратором чата")
)

// In входящие параметры
type In struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// Out результат запроса чатов
type Out struct{}

type DeleteMemberUsecase struct {
	Repo chatt.Repository
}

// DeleteMember удаляет участника чата
func (c *DeleteMemberUsecase) DeleteMember(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectID {
		return Out{}, ErrMemberCannotDeleteHimself
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Subject должен быть главным администратором
	if chat.ChiefID != in.SubjectID {
		return Out{}, ErrSubjectUserIsNotChief
	}

	// Удалить пользователя из чата
	if err = chat.RemoveParticipant(in.UserID); err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	return Out{}, nil
}

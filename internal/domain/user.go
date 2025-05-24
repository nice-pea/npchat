package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// User представляет собой пользователя.
type User struct {
	ID   string
	Name string
	Nick string
}

var (
	ErrUserIDValidate = errors.New("некорректный UUID")
)

// ValidateID проверяет корректность идентификатора пользователя.
func (u User) ValidateID() error {
	if err := uuid.Validate(u.ID); err != nil {
		return errors.Join(err, ErrUserIDValidate)
	}
	return nil // Идентификатор валиден
}

// UserNameMaxLen максимальная длина имени пользователя.
const UserNameMaxLen = 35

// ValidateName проверяет корректность идентификатора пользователя.
func (u User) ValidateName() error {
	// Check if name is empty or contains only whitespace
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name cannot be empty")
	}

	// Проверка на длину имени пользователя
	if len([]rune(u.Name)) > UserNameMaxLen {
		return fmt.Errorf("name cannot be longer than %d characters", UserNameMaxLen)
	}

	// Check for leading or trailing spaces
	if u.Name != strings.TrimSpace(u.Name) {
		return errors.New("name cannot have leading or trailing spaces")
	}

	// Проверка на управляющие символы (кроме обычного пробела)
	for _, r := range u.Name {
		if unicode.IsControl(r) && r != ' ' { // \n, \t, \r и т.д.
			return fmt.Errorf("name cannot contain control characters")
		}
	}

	return nil
}

// UserNickMaxLen максимальная длина ника пользователя.
const UserNickMaxLen = 35

// ValidateNick проверяет корректность идентификатора пользователя.
func (u User) ValidateNick() error {
	// Проверка на пустое имя пользователя
	if u.Nick == "" {
		return nil
	}

	// Проверка на длину имени пользователя
	if len([]rune(u.Nick)) > UserNickMaxLen {
		return fmt.Errorf("nick cannot be longer than %d characters", UserNickMaxLen)
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(u.Nick) == "" {
		return fmt.Errorf("nick cannot consist of whitespace only")
	}

	// Проверка на первый символ
	if !isAllowedLastFirstNickRune(rune(u.Nick[0])) {
		return fmt.Errorf("nick must start with a letter or digit")
	}
	// Проверка на последний символ
	if !isAllowedLastFirstNickRune(rune(u.Nick[len(u.Nick)-1])) {
		return fmt.Errorf("nick must trail with a letter or digit")
	}

	var hasLetters bool
	// Проверка каждого символа в имени пользователя
	for _, r := range u.Nick {
		switch {
		case unicode.IsControl(r):
			return fmt.Errorf("nick cannot contain control characters")
		case unicode.IsSpace(r):
			return fmt.Errorf("nick cannot contain spaces")
		case !isAllowedNickRune(r):
			return fmt.Errorf("nick contains invalid characters")
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}

	if !hasLetters {
		return fmt.Errorf("nick must contain at least one letter or digit")
	}

	return nil // Ник валиден
}

// isAllowedNickRune проверяет, разрешен ли символ в имени пользователя
func isAllowedNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isAllowedNickRune проверяет, разрешен ли символ в имени пользователя
func isAllowedLastFirstNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

type UsersRepository interface {
	// List возвращает список с учетом фильтрации
	List(filter UsersFilter) ([]User, error)

	// Save сохраняет запись
	Save(user User) error

	// Delete удаляет запись
	Delete(id string) error
}

// UsersFilter представляет собой фильтр по пользователям.
type UsersFilter struct {
	ID string // Идентификатор пользователя для фильтрации
}

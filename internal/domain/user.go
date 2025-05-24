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
	ID   string // ID пользователя
	Name string // Имя пользователя
	Nick string // Ник пользователя
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

// ValidateName проверяет корректность имени пользователя.
func (u User) ValidateName() error {
	// Проверить, не является ли имя пустым или содержит только пробелы
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("имя не может быть пустым")
	}

	// Проверка на длину имени пользователя
	if len([]rune(u.Name)) > UserNameMaxLen {
		return fmt.Errorf("длина имени не может превышать %d символов", UserNameMaxLen)
	}

	// Check for leading or trailing spaces
	if u.Name != strings.TrimSpace(u.Name) {
		return errors.New("имя не может содержать начальных или конечных пробелов")
	}

	// Проверка на управляющие символы (кроме обычного пробела)
	for _, r := range u.Name {
		if unicode.IsControl(r) && r != ' ' { // \n, \t, \r и т.д.
			return fmt.Errorf("имя не может содержать управляющих символов")
		}
	}

	return nil
}

// UserNickMaxLen максимальная длина ника пользователя.
const UserNickMaxLen = 35

// ValidateNick проверяет корректность ника пользователя.
func (u User) ValidateNick() error {
	// Проверка на пустое имя пользователя
	if u.Nick == "" {
		return nil
	}

	// Проверка на длину имени пользователя
	if len([]rune(u.Nick)) > UserNickMaxLen {
		return fmt.Errorf("ник не может быть длиннее %d символов", UserNickMaxLen)
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(u.Nick) == "" {
		return fmt.Errorf("ник не может состоять только из пробелов")
	}

	// Проверка на первый символ
	if !isAllowedLastFirstNickRune(rune(u.Nick[0])) {
		return fmt.Errorf("ник должен начинаться с буквы или цифры")
	}
	// Проверка на последний символ
	if !isAllowedLastFirstNickRune(rune(u.Nick[len(u.Nick)-1])) {
		return fmt.Errorf("ник должен заканчиваться буквой или цифрой")
	}

	var hasLetters bool
	// Проверка каждого символа в имени пользователя
	for _, r := range u.Nick {
		switch {
		case unicode.IsControl(r):
			return fmt.Errorf("ник не может содержать управляющие символы")
		case unicode.IsSpace(r):
			return fmt.Errorf("ник не может содержать пробелы")
		case !isAllowedNickRune(r):
			return fmt.Errorf("ник содержит недопустимые символы")
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}

	if !hasLetters {
		return fmt.Errorf("ник должен содержать хотя бы одну букву или цифру")
	}

	return nil // Ник валиден
}

// isAllowedNickRune проверяет, разрешен ли символ в имени пользователя
func isAllowedNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isAllowedLastFirstNickRune проверяет, разрешен ли символ в качестве первого или последнего символа ника
func isAllowedLastFirstNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// UsersRepository интерфейс для работы с репозиторием пользователей.
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
	ID string // ID пользователя для фильтрации
}

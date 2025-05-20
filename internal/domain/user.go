package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

type User struct {
	ID       string
	Name     string
	Username string
}

var (
	ErrUserIDValidate = errors.New("некорректный UUID")
)

func (u User) ValidateID() error {
	if err := uuid.Validate(u.ID); err != nil {
		return errors.Join(err, ErrUserIDValidate)
	}
	return nil
}

const UserNameMaxLen = 35

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

const UserUsernameMaxLen = 35

func (u User) ValidateUsername() error {
	// Проверка на пустое имя пользователя
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Проверка на длину имени пользователя
	if len([]rune(u.Username)) > UserUsernameMaxLen {
		return fmt.Errorf("username cannot be longer than %d characters", UserUsernameMaxLen)
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(u.Username) == "" {
		return fmt.Errorf("username cannot consist of whitespace only")
	}

	// Проверка на первый символ
	if !isAllowedLastFirstUsernameRune(rune(u.Username[0])) {
		return fmt.Errorf("username must start with a letter or digit")
	}
	// Проверка на последний символ
	if !isAllowedLastFirstUsernameRune(rune(u.Username[len(u.Username)-1])) {
		return fmt.Errorf("username must trail with a letter or digit")
	}

	var hasLetters bool
	// Проверка каждого символа в имени пользователя
	for _, r := range u.Username {
		switch {
		case unicode.IsControl(r):
			return fmt.Errorf("username cannot contain control characters")
		case unicode.IsSpace(r):
			return fmt.Errorf("username cannot contain spaces")
		case !isAllowedUsernameRune(r):
			return fmt.Errorf("username contains invalid characters")
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}

	if !hasLetters {
		return fmt.Errorf("username must contain at least one letter or digit")
	}

	return nil
}

// isAllowedUsernameRune проверяет, разрешен ли символ в имени пользователя
func isAllowedUsernameRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isAllowedUsernameRune проверяет, разрешен ли символ в имени пользователя
func isAllowedLastFirstUsernameRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

type UsersRepository interface {
	List(filter UsersFilter) ([]User, error)
	Save(user User) error
	Delete(id string) error
}

type UsersFilter struct {
	ID string
}

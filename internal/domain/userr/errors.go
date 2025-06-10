package userr

import (
	"errors"
	"fmt"
)

var (
	ErrPasswordTooShort       = fmt.Errorf("пароль должен быть не короче %d символов", UserPasswordMinLen)
	ErrPasswordTooLong        = fmt.Errorf("пароль не может быть длиннее %d символов", UserPasswordMaxLen)
	ErrOnlyArabicDigits       = fmt.Errorf("разрешены только арабские цифры (0-9)")
	ErrPasswordInvalidChars   = fmt.Errorf("пароль содержит недопустимые символы")
	ErrPasswordNoUppercase    = fmt.Errorf("пароль должен содержать хотя бы одну заглавную букву")
	ErrPasswordNoLowercase    = fmt.Errorf("пароль должен содержать хотя бы одну строчную букву")
	ErrPasswordNoDigit        = fmt.Errorf("пароль должен содержать хотя бы одну цифру (0-9)")
	ErrLoginTooLong           = fmt.Errorf("логин не может быть длиннее %d символов", UserLoginMaxLen)
	ErrLoginTooShort          = fmt.Errorf("логин не может быть короче %d символов", UserLoginMinLen)
	ErrLoginOnlySpaces        = fmt.Errorf("логин не может состоять только из пробелов")
	ErrLoginStartChar         = fmt.Errorf("логин должен начинаться с буквы или цифры")
	ErrLoginEndChar           = fmt.Errorf("логин должен заканчиваться буквой или цифрой")
	ErrLoginControlChars      = fmt.Errorf("логин не может содержать управляющие символы")
	ErrLoginSpaces            = fmt.Errorf("логин не может содержать пробелы")
	ErrLoginInvalidChars      = fmt.Errorf("логин содержит недопустимые символы")
	ErrLoginNoLetters         = fmt.Errorf("логин должен содержать хотя бы одну букву")
	ErrNameEmpty              = errors.New("имя не может быть пустым")
	ErrNameTooLong            = fmt.Errorf("длина имени не может превышать %d символов", UserNameMaxLen)
	ErrNameSpaces             = errors.New("имя не может содержать начальных или конечных пробелов")
	ErrNameControlChars       = fmt.Errorf("имя не может содержать управляющих символов")
	ErrNickTooLong            = fmt.Errorf("ник не может быть длиннее %d символов", UserNickMaxLen)
	ErrNickOnlySpaces         = fmt.Errorf("ник не может состоять только из пробелов")
	ErrNickStartChar          = fmt.Errorf("ник должен начинаться с буквы или цифры")
	ErrNickEndChar            = fmt.Errorf("ник должен заканчиваться буквой или цифрой")
	ErrNickControlChars       = fmt.Errorf("ник не может содержать управляющие символы")
	ErrNickSpaces             = fmt.Errorf("ник не может содержать пробелы")
	ErrNickInvalidChars       = fmt.Errorf("ник содержит недопустимые символы")
	ErrNickNoLetters          = fmt.Errorf("ник должен содержать хотя бы одну букву или цифру")
	ErrPasswordContainsSpaces = fmt.Errorf("пароль не может содержать пробелы")
)

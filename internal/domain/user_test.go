package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func TestUser_ValidateID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		u := User{ID: ID}
		return u.ValidateID()
	})
}

func Test_UserValidateUsername(t *testing.T) {
	tests := []struct {
		testname string
		username string
		wantErr  bool
	}{
		// Invalid cases
		{testname: "empty username", username: "", wantErr: true},
		{testname: "whitespace only", username: " ", wantErr: true},
		{testname: "control characters", username: "\n\t\r\a\f\v", wantErr: true},
		{testname: "space in middle", username: "first last", wantErr: true},
		{testname: "surrounding spaces", username: " name ", wantErr: true},
		{testname: "digits only", username: "1111", wantErr: true},
		{testname: "special characters", username: "-!@#$%^&*()+}{:\"'?><.", wantErr: true},
		{testname: "trailing underscore", username: "name_", wantErr: true},
		{testname: "leading underscore", username: "_name", wantErr: true},
		{testname: "hyphen in middle", username: "first-last", wantErr: true},
		{testname: "cyrillic letters", username: "ИмяФамилия", wantErr: true},
		{testname: "chinese characters", username: "名字", wantErr: true},
		{testname: "emoji", username: "😊username", wantErr: true}, // или true, если эмодзи запрещены
		{testname: "special unicode (ñ, é)", username: "niño café", wantErr: true},
		{testname: "japanese (kanji + hiragana)", username: "名前なまえ", wantErr: true},
		{testname: "too long name", username: strings.Repeat("a", 35+1), wantErr: true},

		// Valid cases
		{testname: "valid simple name", username: "name", wantErr: false},
		{testname: "valid name with digits", username: "1na1me1", wantErr: false},
		{testname: "valid underscore separated", username: "first_last", wantErr: false},
		{testname: "alphabetic", username: "abcdefghijklmnopqrstuvwxyz", wantErr: false},
		{testname: "alphabetic", username: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			u := User{Username: tt.username}
			if err := u.ValidateUsername(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_UserValidateName(t *testing.T) {
	tests := []struct {
		testname string
		name     string
		wantErr  bool
	}{
		// Invalid cases
		{testname: "empty name", name: "", wantErr: true},
		{testname: "whitespace only", name: " ", wantErr: true},
		{testname: "control characters", name: "q\n\t\r\a\f\vq", wantErr: true},
		{testname: "leading space", name: " Name", wantErr: true},
		{testname: "trailing space", name: "Name ", wantErr: true},
		{testname: "too long name", name: strings.Repeat("名", 35+1), wantErr: true},

		// Valid ASCII cases
		{testname: "single character", name: "q", wantErr: false},
		{testname: "single digit", name: "1", wantErr: false},
		{testname: "lowercase name", name: "name", wantErr: false},
		{testname: "alphanumeric name", name: "Name2", wantErr: false},
		{testname: "name with space", name: "first last", wantErr: false},

		// Unicode cases (valid if supported)
		{testname: "cyrillic letters", name: "ИмяФамилия", wantErr: false},
		{testname: "chinese characters", name: "名字", wantErr: false},
		{testname: "emoji", name: "😊username", wantErr: false},
		{testname: "special unicode (ñ, é)", name: "niño café", wantErr: false},
		{testname: "japanese (kanji + hiragana)", name: "名前なまえ", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			u := User{Name: tt.name}
			if err := u.ValidateName(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

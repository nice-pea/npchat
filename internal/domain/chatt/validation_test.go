package chatt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateChatName(t *testing.T) {
	tests := []struct {
		name     string
		chatName string
		wantErr  bool
	}{
		{
			name:     "пустая строка",
			chatName: "",
			wantErr:  true,
		},
		{
			name:     "превышает лимит в 50 символов",
			chatName: strings.Repeat("a", 51),
			wantErr:  true,
		},
		{
			name:     "содержит пробел в начале",
			chatName: " name",
			wantErr:  true,
		},
		{
			name:     "содержит пробел в конце",
			chatName: "name ",
			wantErr:  true,
		},
		{
			name:     "содержит таб",
			chatName: "na\tme",
			wantErr:  true,
		},
		{
			name:     "содержит новую строку",
			chatName: "na\nme",
			wantErr:  true,
		},
		{
			name:     "содержит только пробелы",
			chatName: " ",
			wantErr:  true,
		},
		{
			name:     "содержит цифры",
			chatName: "1na13me4",
			wantErr:  false,
		},
		{
			name:     "содержит пробел в середине",
			chatName: "na me",
			wantErr:  false,
		},
		{
			name:     "содержит пробелы в середине",
			chatName: "na  me",
			wantErr:  false,
		},
		{
			name:     "содержит знаки",
			chatName: "??na??me.#1432&^$(@",
			wantErr:  false,
		},
		{
			name:     "содержит только знаки",
			chatName: "?>><#(*@$&",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateChatName(tt.chatName); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

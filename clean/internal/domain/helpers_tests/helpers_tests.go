package helpers_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func RunValidateRequiredIDTest(t *testing.T, validate func(string) error) {
	tests := []struct {
		name    string
		ID      string
		wantErr bool
	}{

		{
			name:    "пустая строка как id",
			ID:      "",
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			ID:      "fndsef",
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			ID:      "----",
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			ID:      ",,,,,,,,",
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			ID:      "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			wantErr: true,
		},
		{
			name:    "это uuid",
			ID:      "1cee9c74-a359-425c-b1bb-91c8a35e7b21",
			wantErr: false,
		},
		{
			name:    "это uuid",
			ID:      "0195ba16-f44c-7a3b-b326-94697ec6b00e",
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			ID:      uuid.NewString(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.ID); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

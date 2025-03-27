package domain

import (
	"testing"

	"github.com/google/uuid"
)

func RunValidateAnyIDTest(t *testing.T, validate func(string) error, validateFuncName string) {
	tests := []struct {
		name    string
		anyID   string
		wantErr bool
	}{

		{
			name:    "пустая строка как id",
			anyID:   "",
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			anyID:   "fndsef",
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			anyID:   "----",
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			anyID:   ",,,,,,,,",
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			anyID:   "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			wantErr: true,
		},
		{
			name:    "это uuid",
			anyID:   "1cee9c74-a359-425c-b1bb-91c8a35e7b21",
			wantErr: false,
		},
		{
			name:    "это uuid",
			anyID:   "0195ba16-f44c-7a3b-b326-94697ec6b00e",
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			anyID:   uuid.NewString(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.anyID); (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", validateFuncName, err, tt.wantErr)
			}
		})
	}
}

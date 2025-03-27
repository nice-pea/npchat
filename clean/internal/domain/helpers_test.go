package domain

import (
	"github.com/google/uuid"
)

type fields struct {
	ID string
}
type TestDescription struct {
	name    string
	fields  fields
	wantErr bool
}

func Helper_Test_ValidateID(f func([]TestDescription)) {
	tests := []TestDescription{

		{
			name:    "пустая строка как id",
			fields:  fields{ID: ""},
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			fields:  fields{ID: "fndsef"},
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			fields:  fields{ID: "----"},
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			fields:  fields{ID: ",,,,,,,,"},
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			fields:  fields{ID: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
			wantErr: true,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "1cee9c74-a359-425c-b1bb-91c8a35e7b21"},
			wantErr: false,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "0195ba16-f44c-7a3b-b326-94697ec6b00e"},
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			fields:  fields{ID: uuid.NewString()},
			wantErr: false,
		},
	}
	f(tests)
}

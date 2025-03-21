package domain

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestChat_ValidateID(t *testing.T) {
	type fields struct {
		ID   string
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "пустая строка как id",
			fields:  fields{ID: "", Name: ""},
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			fields:  fields{ID: "fndsef", Name: ""},
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			fields:  fields{ID: "----", Name: ""},
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			fields:  fields{ID: ",,,,,,,,", Name: ""},
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			fields:  fields{ID: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", Name: ""},
			wantErr: true,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "1cee9c74-a359-425c-b1bb-91c8a35e7b21", Name: ""},
			wantErr: false,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "0195ba16-f44c-7a3b-b326-94697ec6b00e", Name: ""},
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			fields:  fields{ID: uuid.NewString(), Name: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chat{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := c.ValidateID(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChat_ValidateName(t *testing.T) {
	type fields struct {
		ID   string
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "пустая строка",
			fields:  fields{Name: ""},
			wantErr: true,
		},
		{
			name:    "превышает лимит в 50 символов",
			fields:  fields{Name: strings.Repeat("a", 51)},
			wantErr: true,
		},
		{
			name:    "содержит пробел в начале",
			fields:  fields{Name: " name"},
			wantErr: true,
		},
		{
			name:    "содержит пробел в конце",
			fields:  fields{Name: "name "},
			wantErr: true,
		},
		{
			name:    "содержит таб",
			fields:  fields{Name: "na\tme"},
			wantErr: true,
		},
		{
			name:    "содержит новую строку",
			fields:  fields{Name: "na\nme"},
			wantErr: true,
		},
		{
			name:    "содержит цифры",
			fields:  fields{Name: "1na13me4"},
			wantErr: false,
		},
		{
			name:    "содержит пробел в середине",
			fields:  fields{Name: "na me"},
			wantErr: false,
		},
		{
			name:    "содержит пробелы в середине",
			fields:  fields{Name: "na  me"},
			wantErr: false,
		},
		{
			name:    "содержит знаки",
			fields:  fields{Name: "??na??me.#1432&^$(@"},
			wantErr: false,
		},
		{
			name:    "содержит только знаки",
			fields:  fields{Name: "?>><#(*@$&"},
			wantErr: false,
		},
		{
			name:    "содержит только пробелы",
			fields:  fields{Name: " "},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chat{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := c.ValidateName(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChatsRepository_List(t *testing.T, newRepository func() ChatsRepository) {
	t.Helper()

	tests := []struct {
		name    string
		init    func(ChatsRepository) ChatsRepository
		filter  ChatsFilter
		wantRes []Chat
		wantErr bool
	}{
		{
			name: "без фильтра в пустом репозитории",
			init: func(repository ChatsRepository) ChatsRepository {
				return repository
			},
			filter:  ChatsFilter{},
			wantRes: []Chat{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newRepository()
			tt.init(repository)
			chats, err := repository.List(tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(chats, tt.wantRes) {
				t.Errorf("List() chats = %v, want %v", chats, tt.filter)
			}
		})
	}
}

func TestChat_ValidateName1(t *testing.T) {
	type fields struct {
		ID   string
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chat{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := c.ValidateName(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

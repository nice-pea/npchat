package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// Uint is an alias for sql.NullInt64 data type
type Uint struct {
	Val uint
	sql sql.NullInt64
}

func (u Uint) IsZero() bool {
	return !u.sql.Valid
}

// Scan implements the Scanner interface for Uint
func (u *Uint) Scan(value interface{}) error {
	var sqlV sql.NullInt64
	if err := sqlV.Scan(value); err != nil {
		return err
	}
	*u = Uint{
		Val: uint(sqlV.Int64),
		sql: sqlV,
	}
	return nil
}

// MarshalJSON for Uint
func (u *Uint) MarshalJSON() ([]byte, error) {
	if !u.sql.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(u.Val)
}

// Value implements the [driver.Valuer] interface.
func (u Uint) Value() (driver.Value, error) {
	return u.sql.Value()
}

// UnmarshalJSON - пользовательская реализация анмаршалинга
func (u *Uint) UnmarshalJSON(data []byte) error {
	var aux *uint // Временная переменная для хранения значения

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux == nil {
		u.sql.Valid = false // Если значение null, устанавливаем Valid в false
		return nil
	}

	u.Val = *aux
	u.sql.Valid = true // Устанавливаем Valid в true, если значение не null
	return nil
}

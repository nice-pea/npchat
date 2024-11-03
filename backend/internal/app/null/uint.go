package null

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

// Uint is an alias for sql.NullInt64 data type
type Uint struct {
	pg pgtype.Uint32
}

// Constructor

func NewUint(val uint, valid bool) Uint {
	return Uint{pg: pgtype.Uint32{
		Uint32: uint32(val),
		Valid:  valid,
	}}
}

// Getters/Setters

func (u *Uint) Val() uint {
	return uint(u.pg.Uint32)
}

func (u *Uint) SetVal(val uint) {
	u.pg.Uint32 = uint32(val)
}

func (u *Uint) Valid() bool {
	return u.pg.Valid
}

func (u *Uint) SetValid(valid bool) {
	u.pg.Valid = valid
}

// Implement pgtype.Uint32 methods

func (u *Uint) ScanUint32(v pgtype.Uint32) error {
	return u.pg.ScanUint32(v)
}

func (u Uint) Uint32Value() (pgtype.Uint32, error) {
	return u.pg.Uint32Value()
}

func (u *Uint) Scan(src any) error {
	return u.pg.Scan(src)
}

func (u Uint) Value() (driver.Value, error) {
	return u.pg.Value()
}

// Implement json

// MarshalJSON for Uint
func (u *Uint) MarshalJSON() ([]byte, error) {
	if !u.pg.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(u.pg.Uint32)
}

// UnmarshalJSON - пользовательская реализация анмаршалинга
func (u *Uint) UnmarshalJSON(data []byte) error {
	var aux *uint32 // Временная переменная для хранения значения

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux == nil {
		u.pg.Valid = false // Если значение null, устанавливаем Valid в false
		return nil
	}

	u.pg.Uint32 = *aux
	u.pg.Valid = true // Устанавливаем Valid в true, если значение не null
	return nil
}

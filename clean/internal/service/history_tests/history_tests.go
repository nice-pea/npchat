package history_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/service"
)

func RunHistoryTest(t *testing.T, newHistory func() service.History) {
	t.Helper()
	type args struct {
		typ   string
		props any
	}
	tests := []struct {
		name      string
		fields    args
		wantPanic bool
	}{
		{
			name:      "nil в props",
			fields:    args{typ: "typ", props: nil},
			wantPanic: false,
		},
		{
			name:      "структура в props",
			fields:    args{typ: "typ", props: struct{ A string }{"A"}},
			wantPanic: false,
		},
		{
			name: "карта со string ключами в props",
			fields: args{
				typ:   "typ",
				props: map[string]any{"A": "321.3"},
			},
			wantPanic: false,
		},
		{
			name: "карта с int ключами в props",
			fields: args{
				typ:   "typ",
				props: map[int]any{321: "A"},
			},
			wantPanic: true,
		},
		{
			name:      "строка в props",
			fields:    args{typ: "typ", props: "some"},
			wantPanic: true,
		},
		{
			name:      "целое число в props",
			fields:    args{typ: "typ", props: 123},
			wantPanic: true,
		},
		{
			name:      "булевый тип в props",
			fields:    args{typ: "typ", props: false},
			wantPanic: true,
		},

		{
			name:      "Slice в props",
			fields:    args{typ: "typ", props: []int{1, 2, 3}},
			wantPanic: true,
		},
		{
			name:      "Int8 в props",
			fields:    args{typ: "typ", props: int8(123)},
			wantPanic: true,
		},
		{
			name:      "Int16 в props",
			fields:    args{typ: "typ", props: int16(521)},
			wantPanic: true,
		},
		{
			name:      "Int32 в props",
			fields:    args{typ: "typ", props: int32(43)},
			wantPanic: true,
		},
		{
			name:      "Int64 в props",
			fields:    args{typ: "typ", props: int64(-53)},
			wantPanic: true,
		},
		{
			name:      "Uint в props",
			fields:    args{typ: "typ", props: uint(644)},
			wantPanic: true,
		},
		{
			name:      "Uint8 в props",
			fields:    args{typ: "typ", props: uint8(3)},
			wantPanic: true,
		},
		{
			name:      "Uint16 в props",
			fields:    args{typ: "typ", props: uint16(65)},
			wantPanic: true,
		},
		{
			name:      "Uint32 в props",
			fields:    args{typ: "typ", props: uint32(76)},
			wantPanic: true,
		},
		{
			name:      "Uint64 в props",
			fields:    args{typ: "typ", props: uint64(76)},
			wantPanic: true,
		},
		{
			name:      "Uintptr в props",
			fields:    args{typ: "typ", props: uintptr(123)},
			wantPanic: true,
		},
		{
			name:      "Float32 в props",
			fields:    args{typ: "typ", props: float32(123.45)},
			wantPanic: true,
		},
		{
			name:      "Float64 в props",
			fields:    args{typ: "typ", props: float64(123.45)},
			wantPanic: true,
		},
		{
			name:      "Complex64 в props",
			fields:    args{typ: "typ", props: complex64(1 + 2i)},
			wantPanic: true,
		},
		{
			name:      "Complex128 в props",
			fields:    args{typ: "typ", props: complex128(1 + 2i)},
			wantPanic: true,
		},
		{
			name:      "Array в props",
			fields:    args{typ: "typ", props: [3]int{1, 2, 3}},
			wantPanic: true,
		},
		{
			name:      "Chan в props",
			fields:    args{typ: "typ", props: make(chan int)},
			wantPanic: true,
		},
		{
			name:      "Func в props",
			fields:    args{typ: "typ", props: func() {}},
			wantPanic: true,
		},
		{
			name:      "Pointer в props",
			fields:    args{typ: "typ", props: new(int)},
			wantPanic: true,
		},
		{
			name:      "UnsafePointer в props",
			fields:    args{typ: "typ", props: new(int)},
			wantPanic: true,
		},
		{
			name:      "не пустая строка в typ",
			fields:    args{typ: "3aRc,.3%@)(#&$"},
			wantPanic: false,
		},
		{
			name:      "пустая строка в typ",
			fields:    args{typ: ""},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHistory()
			if tt.wantPanic {
				assert.Panics(t, func() {
					h.Write(tt.fields.typ, tt.fields.props)
				})
			} else {
				assert.NotPanics(t, func() {
					h.Write(tt.fields.typ, tt.fields.props)
				})
			}
		})
	}
}

package extend

import "testing"

func TestExtend_ResolveConflicts(t *testing.T) {
	type testCase[T any] struct {
		name    string
		e       Params[T]
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "ok",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"b", "c"}, Fn: nil},
					{Key: "c", Deps: []string{"b"}, Fn: nil},
					{Key: "b", Deps: []string{"f"}, Fn: nil},
					{Key: "f", Deps: nil, Fn: nil},
				},
			},
			wantErr: false,
		},
		{
			name: "ErrCantCreateDep",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"b", "c"}, Fn: nil},
					{Key: "c", Deps: []string{"b"}, Fn: nil},
					{Key: "b", Deps: []string{"f"}, Fn: nil},
					{Key: "f", Deps: []string{"x"}, Fn: nil},
				},
			},
			wantErr: true,
		},
		{
			name: "Collision",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"x"}, Fn: nil},
					{Key: "x", Deps: []string{"a"}, Fn: nil},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.e.sort(); (err != nil) != tt.wantErr {
				t.Errorf("sort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
					{Key: "a", Deps: []string{"b", "c"}},
					{Key: "c", Deps: []string{"b"}},
					{Key: "b", Deps: []string{"f"}},
					{Key: "f", Deps: nil},
				},
			},
			wantErr: false,
		},
		{
			name: "ok",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"b", "c"}},
					{Key: "f", Deps: nil},
					{Key: "b", Deps: []string{"f"}},
					{Key: "c", Deps: []string{"b"}},
				},
			},
			wantErr: false,
		},
		{
			name: "ok",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "f", Deps: nil},
					{Key: "c", Deps: []string{"b"}},
					{Key: "b", Deps: []string{"f"}},
					{Key: "a", Deps: []string{"b", "c"}},
				},
			},
			wantErr: false,
		},
		{
			name: "ErrCantCreateDep",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"b", "c"}},
					{Key: "c", Deps: []string{"b"}},
					{Key: "b", Deps: []string{"f"}},
					{Key: "f", Deps: []string{"x"}},
				},
			},
			wantErr: true,
		},
		{
			name: "Collision",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"x"}},
					{Key: "x", Deps: []string{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name:    "Collision",
			e:       Params[int]{},
			wantErr: false,
		},
		{
			name: "Collision",
			e: Params[int]{
				Fields: []Field[int]{
					{Key: "a", Deps: []string{"x"}},
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

package extend

import "testing"

func TestExtend_ResolveConflicts(t *testing.T) {
	type testCase struct {
		name    string
		e       Params
		wantErr bool
	}
	tests := []testCase{
		{
			name: "ok",
			e: Params{
				Fields: []Field{
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
			e: Params{
				Fields: []Field{
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
			e: Params{
				Fields: []Field{
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
			e: Params{
				Fields: []Field{
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
			e: Params{
				Fields: []Field{
					{Key: "a", Deps: []string{"x"}},
					{Key: "x", Deps: []string{"a"}},
				},
			},
			wantErr: true,
		},
		{
			name:    "Collision",
			e:       Params{},
			wantErr: false,
		},
		{
			name: "Collision",
			e: Params{
				Fields: []Field{
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

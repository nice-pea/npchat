package extend

import (
	"errors"
	"slices"
)

type Field struct {
	Key  string
	Deps []string
	Fn   func() error
}

type Params struct {
	Fields []Field
}

func (p Params) Run() error {
	if sorted, err := p.sort(); err != nil {
		return err
	} else {
		for _, field := range sorted {
			if err = field.Fn(); err != nil {
				return err
			}
		}
	}

	return nil
}

var ErrCantCreateDep = errors.New("can't create dep")

func (p Params) sort() ([]Field, error) {
	exchange := make(map[string]Field, len(p.Fields))
	for _, field := range p.Fields {
		exchange[field.Key] = field
	}

	created := make([]Field, 0, len(p.Fields))
	for i := 0; i < len(exchange) && len(exchange) != 0; i++ {
		appended := false
		for _, field := range exchange {
			if !depsAlreadyCreated(field, created) {
				continue
			}
			appended = true
			created = append(created, field)
			delete(exchange, field.Key)
		}
		if !appended {
			return nil, ErrCantCreateDep
		}
	}

	return created, nil
}

func depsAlreadyCreated(field Field, created []Field) bool {
	createdKeys := make([]string, len(created))
	for i, createdField := range created {
		createdKeys[i] = createdField.Key
	}
	for _, dep := range field.Deps {
		if !slices.Contains(createdKeys, dep) {
			return false
		}
	}

	return true
}

//func (p Params) requiredBys() map[string]map[string]struct{} {
//	requiredBys := make(map[string]map[string]struct{}, len(p.Fields))
//	for _, field := range p.Fields {
//		rb := make(map[string]struct{}, len(field.Deps))
//		for _, dep := range field.Deps {
//			rb[dep] = struct{}{}
//		}
//		requiredBys[field.Key] = rb
//	}
//
//	return requiredBys
//}

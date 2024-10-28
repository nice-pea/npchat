package extend

import (
	"errors"
	"slices"
)

type Field[T any] struct {
	Key  string
	Deps []string
	Fn   func(state T) error
}

type Params[T any] struct {
	State  T
	Fields []Field[T]
}

func (p Params[T]) Run() error {
	if sorted, err := p.sort(); err != nil {
		return err
	} else {
		for _, field := range sorted {
			if err = field.Fn(p.State); err != nil {
				return err
			}
		}
	}

	return nil
}

var ErrCantCreateDep = errors.New("can't create dep")

func (p Params[T]) sort() ([]Field[T], error) {
	exchange := make(map[string]Field[T], len(p.Fields))
	for _, field := range p.Fields {
		exchange[field.Key] = field
	}

	created := make([]Field[T], 0, len(p.Fields))
	i := -1
	strike := 0
	appended := false
	for len(exchange) != 0 {
		i++
		for _, field := range exchange {
			if !depsAlreadyCreated(field, created) {
				continue
			}
			appended = true
			created = append(created, field)
			delete(exchange, field.Key)
		}
		// По всем полям прошел цикл, а значит хоть одну зависимость добавил.
		// Есть недоработка, в зависимости от порядка могут быть ложные срабатывания.
		// TODO: fix
		// TODO: а возможно и нет
		if i+1 >= len(p.Fields) {
			if appended {
				strike = 0
			} else if strike >= 2 {
				return nil, ErrCantCreateDep
			} else {
				strike += 1
			}
			i = -1
		}
	}

	return created, nil
}

func depsAlreadyCreated[T any](field Field[T], created []Field[T]) bool {
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

//func (p Params[T]) requiredBys() map[string]map[string]struct{} {
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

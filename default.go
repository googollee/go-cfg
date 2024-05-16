package cfg

import (
	"context"
	"fmt"
	"reflect"
)

type defaultSource struct {
	values map[string]nullValue
}

func fromDefault() Source {
	return &defaultSource{
		values: make(map[string]nullValue),
	}
}

func (s *defaultSource) Setup(t reflect.Type) error {
	v := reflect.New(t)

	return walkFields(v, "", ".", nil, func(meta fieldMeta, v reflect.Value) error {
		nv := newNullValue(v.Kind())

		if meta.Default != "" {
			if err := nv.UnmarshalText([]byte(meta.Default)); err != nil {
				return fmt.Errorf("can't parse default value %q for field %q: %w", meta.Default, meta.FullKey, err)
			}
		}

		s.values[meta.FullKey] = nv

		return nil
	})
}
func (s *defaultSource) Parse(ctx context.Context, v any) error {
	return walkFields(reflect.ValueOf(v), "", ".", nil, func(meta fieldMeta, r reflect.Value) error {
		defaultValue, ok := s.values[meta.FullKey]
		if !ok {
			return nil
		}

		defaultValue.CopyTo(r)

		return nil
	})
}

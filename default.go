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
		var fv nullValue
		switch v.Kind() {
		case reflect.Int:
			fv = &nullInt[int]{}
		case reflect.Int64:
			fv = &nullInt[int64]{}
		case reflect.String:
			fv = &nullString{}
		}

		if meta.Default != "" {
			if err := fv.UnmarshalText([]byte(meta.Default)); err != nil {
				return fmt.Errorf("can't parse default value for key %q: %w", meta.Key, err)
			}
		}

		s.values[meta.FullKey] = fv

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

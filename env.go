package cfg

import (
	"context"
	"os"
	"reflect"
	"strings"
)

type EnvOption func(*envSource)

func EnvSplitter(splitter string) EnvOption {
	return func(src *envSource) {
		src.splitter = splitter
	}
}

type envSource struct {
	prefix   string
	splitter string
}

func FromEnv(prefix string, opt ...EnvOption) Source {
	ret := &envSource{
		prefix:   prefix,
		splitter: "_",
	}

	for _, f := range opt {
		f(ret)
	}

	return ret
}

func (s *envSource) Setup(t reflect.Type) error {
	return nil
}

func (s *envSource) Parse(ctx context.Context, v any) error {
	return s.parseValue(reflect.ValueOf(v), s.prefix+s.splitter)
}

func (s *envSource) parseValue(v reflect.Value, prefix string) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := prefix + field.Tag.Get("cfg")
		key = strings.ToUpper(key)

		if field.Type.Kind() == reflect.Struct {
			if err := s.parseValue(v.Field(i), key+s.splitter); err != nil {
				return err
			}
			continue
		}

		value, ok := os.LookupEnv(key)
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			var got int
			if err := converToInt(value, &got); err != nil {
				return err
			}
			v.Field(i).Set(reflect.ValueOf(got))
		case reflect.Int64:
			var got int64
			if err := converToInt(value, &got); err != nil {
				return err
			}
			v.Field(i).Set(reflect.ValueOf(got))
		case reflect.String:
			var got string
			got = value
			v.Field(i).Set(reflect.ValueOf(got))
		}
	}

	return nil
}

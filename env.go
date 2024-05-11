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
	return walkFields(reflect.ValueOf(v), s.prefix, s.splitter, nil, func(key string, index []int, v reflect.Value) error {
		key = strings.ToUpper(key)

		value, ok := os.LookupEnv(key)
		if !ok {
			return nil
		}

		v = digPtr(v)

		switch v.Kind() {
		case reflect.Int:
			var got int
			if err := converToInt(value, &got); err != nil {
				return err
			}
			v.Set(reflect.ValueOf(got))
		case reflect.Int64:
			var got int64
			if err := converToInt(value, &got); err != nil {
				return err
			}
			v.Set(reflect.ValueOf(got))
		case reflect.String:
			var got string
			got = value
			v.Set(reflect.ValueOf(got))
		}

		return nil
	})
}

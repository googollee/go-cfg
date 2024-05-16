package cfg

import (
	"context"
	"fmt"
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
	return walkFields(reflect.ValueOf(v), s.prefix, s.splitter, nil, func(meta fieldMeta, v reflect.Value) error {
		key := strings.ToUpper(meta.FullKey)

		value, ok := os.LookupEnv(key)
		if !ok {
			return nil
		}

		v = digPtr(v)

		nv := newNullValue(v.Type())
		if err := nv.UnmarshalText([]byte(value)); err != nil {
			return fmt.Errorf("can't parse value %q of env %q: %w", value, key, err)
		}

		nv.CopyTo(v)

		return nil
	})
}

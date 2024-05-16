package cfg

import (
	"context"
	"flag"
	"fmt"
	"reflect"
)

type FlagOption func(*flagSource)

func FlagWithFlagSet(set FlagSet) FlagOption {
	if set == nil {
		set = flag.CommandLine
	}

	return func(s *flagSource) {
		s.flagset = set
	}
}

func FlagSplitter(splitter string) FlagOption {
	return func(s *flagSource) {
		s.splitter = splitter
	}
}

type flagSource struct {
	prefix   string
	splitter string
	flagset  FlagSet
	values   []nullValue
}

func FromFlag(prefix string, opt ...FlagOption) Source {
	ret := &flagSource{
		prefix:   prefix,
		splitter: ".",
		flagset:  flag.CommandLine,
	}

	for _, o := range opt {
		o(ret)
	}

	return ret
}

func (s *flagSource) Setup(t reflect.Type) error {
	return walkFields(reflect.New(t), s.prefix, s.splitter, nil, func(meta fieldMeta, v reflect.Value) error {
		v = digPtr(v)

		nv := newNullValue(v.Type())

		if meta.Default != "" {
			if err := nv.UnmarshalText([]byte(meta.Default)); err != nil {
				return fmt.Errorf("can't parse default value %q for flag %q: %w", meta.Default, meta.FullKey, err)
			}
		}

		s.flagset.TextVar(nv, meta.FullKey, nv, meta.Usage)
		nv.Init(meta.Index)
		s.values = append(s.values, nv)

		return nil
	})
}

func (s *flagSource) Parse(ctx context.Context, v any) error {
	value := digPtr(reflect.ValueOf(v))

	for _, flagValue := range s.values {
		fieldValue := value.FieldByIndex(flagValue.Index())

		if !flagValue.Valid() {
			continue
		}

		flagValue.CopyTo(fieldValue)
	}

	return nil
}

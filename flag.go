package cfg

import (
	"context"
	"flag"
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
	values   []flagValue
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
	return walkFields(reflect.New(t), s.prefix, s.splitter, nil, func(key string, index []int, v reflect.Value) error {
		v = digPtr(v)

		switch v.Kind() {
		case reflect.Int:
			nv := &nullInt[int]{}
			nv.index = index
			s.flagset.TextVar(nv, key, nv, "")
			s.values = append(s.values, nv)
		case reflect.Int64:
			nv := &nullInt[int64]{}
			nv.index = index
			s.flagset.TextVar(nv, key, nv, "")
			s.values = append(s.values, nv)
		case reflect.String:
			nv := &nullString{}
			nv.index = index
			s.flagset.TextVar(nv, key, nv, "")
			s.values = append(s.values, nv)
		default:
			panic("unknown type " + v.Type().String())
		}
		return nil
	})
}

func (s *flagSource) Parse(ctx context.Context, v any) error {
	value := digPtr(reflect.ValueOf(v))

	for _, fv := range s.values {
		if !fv.Valid() {
			continue
		}

		fv.CopyTo(value.FieldByIndex(fv.Index()))
	}

	return nil
}

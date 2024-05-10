package cfg

import (
	"context"
	"flag"
	"fmt"
	"reflect"
	"slices"
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
	return s.setupType(t, nil, s.prefix)
}

func (s *flagSource) Parse(ctx context.Context, v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for _, fv := range s.values {
		fmt.Println(fv)
		if !fv.Valid() {
			continue
		}

		fv.CopyTo(value.FieldByIndex(fv.Index()))
	}

	return nil
}

func (s *flagSource) setupType(t reflect.Type, index []int, prefix string) error {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		key := prefix + s.splitter + field.Tag.Get("cfg")
		switch field.Type.Kind() {
		case reflect.Struct:
			if err := s.setupType(field.Type, slices.Clip(append(index, i)), key); err != nil {
				return err
			}
		case reflect.Int:
			nv := &nullInt[int]{}
			s.flagset.TextVar(nv, key, &nullInt[int]{}, "")
			nv.index = append(index, i)
			s.values = append(s.values, nv)
		case reflect.Int64:
			nv := &nullInt[int64]{}
			s.flagset.TextVar(nv, key, &nullInt[int]{}, "")
			nv.index = append(index, i)
			s.values = append(s.values, nv)
		case reflect.String:
			nv := &nullString{}
			s.flagset.TextVar(nv, key, &nullString{}, "")
			nv.index = append(index, i)
			s.values = append(s.values, nv)
		default:
			panic("unknown type " + field.Type.String())
		}
	}

	return nil
}

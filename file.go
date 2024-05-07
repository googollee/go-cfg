package cfg

import (
	"flag"
	"io"
	"reflect"
)

type FileFormat interface {
	ExtNames() []string
	MimeNames() []string
	TagName() string
	Parse(r io.Reader, a any) error
}

type FileOption func(s *fileSource) error

func FileFlagSet(set FlagSet) FileOption {
	if set == nil {
		set = flag.CommandLine
	}

	return func(s *fileSource) error {
		s.flagset = set
		return nil
	}
}

type fileSource struct {
	flagset FlagSet
}

func FromFile(flagName, flagValue, flagUsage string, opt ...FileOption) (Source, error) {
	ret := &fileSource{
		flagset: flag.CommandLine,
	}

	for _, o := range opt {
		if err := o(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (s *fileSource) Setup(t reflect.Type) error {
	return nil
}

func (s *fileSource) Parse(v any) error {
	return nil
}

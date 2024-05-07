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

func FileSplitter(splitter string) FileOption {
	return func(s *fileSource) error {
		s.splitter = splitter
		return nil
	}
}

type fileSource struct {
	prefix   string
	splitter string
	flagset  FlagSet
	values   any
	keys     [][]string
}

func FromFlagFile(prefix, flagName, flagValue, flagUsage string, opt ...FileOption) (Source, error) {
	ret := &fileSource{
		prefix:   prefix,
		splitter: ".",
		flagset:  flag.CommandLine,
	}

	for _, o := range opt {
		if err := o(ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func (s *fileSource) Setup(t reflect.Type) error {
	s.values = reflect.New(t)
	return nil
}

func (s *fileSource) Parse(v any) error {
	return nil
}

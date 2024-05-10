package cfg

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type FileDecoder interface {
	ExtNames() []string
	MimeNames() []string
	TagName() string
	Decode(r io.Reader, a any) error
}

type JSON struct{}

func (JSON) ExtNames() []string  { return []string{"json", "js"} }
func (JSON) MimeNames() []string { return []string{"application/javascript", "application/json"} }
func (JSON) TagName() string     { return "json" }

func (JSON) Decode(r io.Reader, a any) error {
	return json.NewDecoder(r).Decode(a)
}

type FileOption func(s *fileSource)

func FileFlagSet(set FlagSet) FileOption {
	return func(src *fileSource) {
		src.flagset = set
	}
}

func FileDecoders(decoder ...FileDecoder) FileOption {
	return func(src *fileSource) {
		src.decoders = decoder
	}
}

type fileSource struct {
	flagset  FlagSet
	decoders []FileDecoder

	filename    string
	ext2decoder map[string]FileDecoder
}

func FromFile(flagName, flagValue, flagUsage string, opt ...FileOption) Source {
	ret := &fileSource{
		flagset:     flag.CommandLine,
		decoders:    []FileDecoder{JSON{}},
		ext2decoder: make(map[string]FileDecoder),
	}

	for _, optFn := range opt {
		optFn(ret)
	}

	for _, decoder := range ret.decoders {
		for _, ext := range decoder.ExtNames() {
			ret.ext2decoder[ext] = decoder
		}
	}

	ret.flagset.StringVar(&ret.filename, flagName, flagValue, flagUsage)

	return ret
}

func (s *fileSource) Setup(t reflect.Type) error {
	return nil
}

func (s *fileSource) Parse(ctx context.Context, v any) error {
	ext := strings.TrimLeft(filepath.Ext(s.filename), ".")
	decoder := s.ext2decoder[ext]
	if decoder == nil {
		return fmt.Errorf("parse %q error: no decoder", s.filename)
	}

	f, err := os.Open(s.filename)
	if err != nil {
		return fmt.Errorf("parse %q error: %w", s.filename, err)
	}
	defer f.Close()

	return decoder.Decode(f, v)
}

package cfg

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlagSource(t *testing.T) {
	type Config struct {
		Str   string `cfg:"str"`
		I     int    `cfg:"i"`
		Inner struct {
			Str string `cfg:"str"`
			I   int    `cfg:"i"`
		} `cfg:"inner"`
	}
	set := flag.NewFlagSet("", flag.PanicOnError)
	src := FromFlag("flag", FlagSplitter("."), FlagWithFlagSet(set))

	if err := src.Setup(reflect.TypeOf(Config{})); err != nil {
		t.Fatal("setup error:", err)
	}

	set.Parse([]string{"--flag.str", "out_str", "--flag.i", "10", "--flag.inner.str", "inner_str"})
	var config Config
	config.Inner.I = 20
	if err := src.Parse(context.Background(), &config); err != nil {
		t.Fatal("parse error:", err)
	}

	want := Config{
		Str: "out_str",
		I:   10,
		Inner: struct {
			Str string "cfg:\"str\""
			I   int    "cfg:\"i\""
		}{
			Str: "inner_str",
			I:   20,
		},
	}

	if diff := cmp.Diff(config, want); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

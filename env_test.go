package cfg

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvSource(t *testing.T) {
	type Config struct {
		Str   string `cfg:"str"`
		I     int    `cfg:"i"`
		Inner struct {
			Str string `cfg:"str"`
			I   int    `cfg:"i"`
		} `cfg:"inner"`
	}

	for key, value := range map[string]string{
		"TEST_STR":       "outer_str",
		"TEST_I":         "10",
		"TEST_INNER_STR": "inner_str",
	} {
		os.Setenv(key, value)
	}

	src := FromEnv("test", EnvSplitter("_"))
	var config Config
	config.Inner.I = 20
	if err := src.Parse(context.Background(), &config); err != nil {
		t.Fatal("parse error:", err)
	}

	want := Config{
		Str: "outer_str",
		I:   10,
		Inner: struct {
			Str string `cfg:"str"`
			I   int    `cfg:"i"`
		}{
			Str: "inner_str",
			I:   20,
		},
	}

	if diff := cmp.Diff(config, want); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

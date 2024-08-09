package cfg

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvSource(t *testing.T) {
	type Inner struct {
		Str string `cfg:"str"`
		I   int    `cfg:"i"`
	}
	type Config struct {
		Str   string `cfg:"str"`
		I     int    `cfg:"i"`
		Inner Inner  `cfg:"inner"`
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
	config.I = 20       // Env should overwrite this field
	config.Inner.I = 20 // Env doesn't exist, should keep the value of this field
	if err := src.Parse(context.Background(), &config); err != nil {
		t.Fatal("parse error:", err)
	}

	want := Config{
		Str: "outer_str",
		I:   10,
		Inner: Inner{
			Str: "inner_str",
			I:   20,
		},
	}

	if diff := cmp.Diff(config, want); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

func TestEnvParseField(t *testing.T) {
	tests := []struct {
		input string
		want  any
	}{
		{"str", "str"},
		{"120", 120},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			meta := fieldMeta{
				FullKey: "TEST_PARSE_FIELD",
			}
			os.Setenv(meta.FullKey, tc.input)

			src := envSource{}
			got := reflect.New(reflect.TypeOf(tc.want)).Elem()

			if err := src.parseField(meta, got); err != nil {
				t.Fatal("src.parseField() error:", err)
			}
			if !reflect.DeepEqual(got.Interface(), tc.want) {
				t.Errorf("src.parseField() = %v, want: %v", got.Interface(), tc.want)
			}
		})
	}
}

func TestEnvParseFieldError(t *testing.T) {
	tests := []struct {
		input string
		want  any
		err   string
	}{
		{"str", 120, `can't parse value "str" of env "TEST_PARSE_FIELD": convert "str" to int error: strconv.ParseInt: parsing "str": invalid syntax`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			meta := fieldMeta{
				FullKey: "TEST_PARSE_FIELD",
			}
			os.Setenv(meta.FullKey, tc.input)

			src := envSource{}
			got := reflect.New(reflect.TypeOf(tc.want)).Elem()

			err := src.parseField(meta, got)
			if err == nil {
				t.Fatal("src.parseField() should return an error, but not")
			}
			if diff := cmp.Diff(tc.err, err.Error()); diff != "" {
				t.Errorf("src.parseField() error diff:\n%s", diff)
			}
		})
	}
}

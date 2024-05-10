package cfg

import (
	"context"
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFileSource(t *testing.T) {
	type Config struct {
		Str string `json:"str"`
		I   int    `json:"i"`
	}

	set := flag.NewFlagSet("", flag.ContinueOnError)

	src := FromFile("config", "./config.json", "Config filepath", FileFlagSet(set), FileDecoders(JSON{}))
	if err := set.Parse([]string{"--config", "./testdata/config.json"}); err != nil {
		t.Fatal("flag set parses error:", err)
	}

	var config Config
	if err := src.Parse(context.Background(), &config); err != nil {
		t.Fatal("parse json config error:", err)
	}

	want := Config{
		Str: "string",
		I:   10,
	}
	if diff := cmp.Diff(config, want); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

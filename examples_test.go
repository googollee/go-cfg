package cfg_test

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/googollee/go-cfg"
)

func TestLoadJSON(t *testing.T) {
	type Config struct {
		Key1  string `cfg:"key1"`
		Key2  int    `cfg:"key2"`
		Embed struct {
			Key string `cfg:"key"`
		} `cfg:"inner"`
	}

	strConfig := `{'key1':'value1', 'key2':10, 'inner':{'key': 'value'}}`

	var config Config
	if err := cfg.Load(strings.NewReader(strConfig), &config, cfg.F("yaml", yaml.NewDecoder), cfg.F("json", json.NewDecoder)); err != nil {
		panic(err)
	}

	if got, want := config.Key1, "value1"; got != want {
		t.Errorf("config.Key1 = %q, want: %q", got, want)
	}
	if got, want := config.Key2, 10; got != want {
		t.Errorf("config.Key2 = %v, want: %v", got, want)
	}
	if got, want := config.Embed.Key, "value"; got != want {
		t.Errorf("config.Embed.Key = %q, want: %q", got, want)
	}
}

func TestLoadYAML(t *testing.T) {
	type Config struct {
		Key1  string `cfg:"key1"`
		Key2  int    `cfg:"key2"`
		Embed struct {
			Key string `cfg:"key"`
		} `cfg:"inner"`
	}

	strConfig := `key1: value1
key2: 10
inner:
  key: value
`

	var config Config
	if err := cfg.Load(strings.NewReader(strConfig), &config, cfg.F("yaml", yaml.NewDecoder), cfg.F("json", json.NewDecoder)); err != nil {
		panic(err)
	}

	if got, want := config.Key1, "value1"; got != want {
		t.Errorf("config.Key1 = %q, want: %q", got, want)
	}
	if got, want := config.Key2, 10; got != want {
		t.Errorf("config.Key2 = %v, want: %v", got, want)
	}
	if got, want := config.Embed.Key, "value"; got != want {
		t.Errorf("config.Embed.Key = %q, want: %q", got, want)
	}
}

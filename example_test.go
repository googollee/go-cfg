package cfg_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/googollee/go-cfg"
)

type Config struct {
	Str     string `cfg:"str,,string value from file"`
	I       int    `cfg:"int,,int value from file"`
	FromEnv struct {
		Str string `cfg:"str,none,string value from env"`
	} `cfg:"from_env"`
	FromFlag struct {
		Str string `cfg:"str,none,string value from flags"`
	} `cfg:"from_flag"`
	WithDefault struct {
		Str string `cfg:"str,default_value,string value with default"`
	} `cfg:"with_default"`
}

// An example to load `Config` from sources orderly:
// - from a config file
// - from env variables
// - from flags
//
// The same field in the flag overwrites the env variable, and the same file in the env variable overwrites the value in the config file.
func ExampleParser() {
	os.Setenv("DEMO_FROM_ENV_STR", "string_from_env")

	set := flag.NewFlagSet("demo", flag.PanicOnError)

	parser := cfg.Parse[Config](
		cfg.FromFile("config", "./testdata/config.json", "config file",
			cfg.FileDecoders(cfg.JSON{}),
			cfg.FileFlagSet(set)),
		cfg.FromEnv("DEMO", cfg.EnvSplitter("_")),
		cfg.FromFlag("demo", cfg.FlagSplitter("."), cfg.FlagWithFlagSet(set)),
	)

	if err := set.Parse([]string{
		"--config", "./testdata/config.json",
		"--demo.from_flag.str", "string_from_flag",
	}); err != nil {
		fmt.Println("flag error:", err)
		return
	}

	var buf bytes.Buffer
	set.SetOutput(&buf)
	set.Usage()
	fmt.Println(strings.ReplaceAll(buf.String(), "\t", "  "))

	var config Config
	if err := parser.Parse(context.Background(), &config); err != nil {
		fmt.Println("load config error:", err)
		return
	}

	fmt.Println("config:", config)

	// Output:
	// Usage of demo:
	//   -config string
	//       config file (default "./testdata/config.json")
	//   -demo.from_env.str value
	//       string value from env (default none)
	//   -demo.from_flag.str value
	//       string value from flags (default none)
	//   -demo.int value
	//       int value from file
	//   -demo.str value
	//       string value from file
	//   -demo.with_default.str value
	//       string value with default (default default_value)
	//
	// config: {string 10 {string_from_env} {string_from_flag} {default_value}}
}

package cfg_test

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/googollee/go-cfg"
)

type Config struct {
	Str   string `cfg:"str"`
	I     int    `cfg:"i"`
	Inner struct {
		Str string `cfg:"str"`
		I   int    `cfg:"i"`
	} `cfg:"inner"`
}

// An example to load `Config` from sources orderly:
// - from a config file
// - from env variables
// - from flags
//
// The same field in the flag overwrites the env variable, and the same file in the env variable overwrites the value in the config file.
func ExampleParser() {
	os.Setenv("DEMO_INNER_I", "20")
	os.Setenv("DEMO_INNER_STR", "inner_str")

	set := flag.NewFlagSet("demo", flag.PanicOnError)

	parser := cfg.Parse[Config](
		cfg.FromFile("config", "./testdata/config.json", "config file",
			cfg.FileDecoders(cfg.JSON{}),
			cfg.FileFlagSet(set)),
		cfg.FromEnv("DEMO", cfg.EnvSplitter("_")),
		cfg.FromFlag("demo", cfg.FlagSplitter("."), cfg.FlagWithFlagSet(set)),
	)

	help := set.Bool("help", false, "show help")
	if err := set.Parse([]string{
		"--config", "./testdata/config.json",
		"--demo.inner.str", "overwrited_inner_str",
		"--demo.str", "overwrited_out_str",
	}); err != nil {
		fmt.Println("flag error:", err)
		return
	}

	if *help {
		flag.Usage()
		return
	}

	var config Config
	if err := parser.Parse(context.Background(), &config); err != nil {
		fmt.Println("load config error:", err)
		return
	}

	fmt.Println("config:", config)

	// Output:
	// config: {overwrited_out_str 10 {overwrited_inner_str 20}}
}

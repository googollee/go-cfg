package cfg_test

import (
	"flag"
	"fmt"

	"github.com/googollee/go-cfg"
)

type Config struct {
	Str string `cfg:"str"`
	Int int    `cfg:"int"`
}

// An example to load `Config` from sources orderly:
// - from a config file
// - from env variables
// - from flags
//
// The same field in the flag overwrites the env variable, and the same file in the env variable overwrites the value in the config file.
func ExampleParser() {
	parser := cfg.Parse[Config](
		cfg.FromFlags("demo", cfg.FlagSplitter(".")),
		cfg.FromEnv("DEMO", cfg.EnvSplitter("_")),
		cfg.FromFile("config", "./config.json", "config file",
			cfg.FileFormat(cfg.JSON, cfg.YAML, cfg.TOML),
			cfg.FileFlag(flag.CommandLine)),
	)

	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var config Config
	if err := parser.Parse(&config); err != nil {
		fmt.Println("load config error:", err)
		return
	}

	fmt.Println("config:", config)
}

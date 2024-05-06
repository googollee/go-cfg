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
func ExampleConfig() {
	fromFlags := cfg.FromFlags[Config]("demo_")

	cfgFile := flag.String("config", "./config.yaml", "config file")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var config Config
	if err := cfg.Parse(&config, cfg.FromFile(*cfgFile /* optional file formats: cfg.YAML, cfg.JSON, cfg.TOML */), cfg.FromEnv("DEMO_"), fromFlags); err != nil {
		fmt.Println("parse config file %q error: %v", *cfgFile, err)
		return
	}

	fmt.Println("config:", config)
}

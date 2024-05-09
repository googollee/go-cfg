//go:build demo
package cfg_test

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/googollee/go-cfg"
)

type DB struct {
	Addr string
}

func (db *DB) Schema(schema *cfg.Schema) {
	schema.NameConvertor(cfg.Gonic)
	schema.Description("Database configuration")
	cfg.Field(schema, &db.Addr).
		Description("The address of database").
		Pattern(regexp.MustCompile("mysql://.*"))
}

type Config struct {
	Str     string `cfg:"str"`
	Int     int    `cfg:"int"`
	Account struct {
		Email string `cfg:"email"`
	} `cfg:"account"`
	DB DB
}

func (c *Config) Schema(schema *cfg.Schema) {
	cfg.Field(schema, &c.Str).
		Default("default").
		Description("Some string value")
	cfg.Field(schema, &c.Account).
		Description("The account infomation")
	cfg.Field(schema, &c.Account.Email).
		Description("The email of the account").
		Pattern(regexp.MustCompile("[a-zA-Z0-9.-_]*@[a-zA-Z0-9.-_]*"))
	cfg.Field(schema, &c.DB).
		Name("db")
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

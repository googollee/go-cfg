/*
package cfg implements a framework to load/parse configuration from a file, environment or flags.

Usage:

  - Register functions ([RegisterInitializer] and [RegisterValue]) should be called in `init()` of a package, before [Init] calls.
  - [Init] function should be called at the beginning of `main()`, before calling functions in other sub-packages.
  - [Init] function must not be called in `init()`, because other sub-packages may not be initialized at that time.
*/
package cfg

import "context"

// Initializer is the interface that wraps a instance that provides configuration information and saves the value provided by the configuration.
//
// The `cfg` library parses values from a file, environment or flags and stores into registered instances. If the instance is an Initializer, the method `Initializer.Init()` is called after storing values. The Initializer provider could initialize global instances in this method.
type Initializer interface {
	Init(ctx context.Context) error
}

/*
RegisterInitializer registers an Initializer instance `value` with the `name` as the scope name. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`. The `value` provider could initialize global instances in `value.Init()`.

Example:

	package subpack

	var singletonInstance int

	type config struct {
		Int int `cfg:"int"`
	}

	func (c config) Init(ctx context.Context) error {
		singletonInstance = c.Int
	}

	func init() {
		cfg.RegisterInitializer[config]("subpack")
	}
*/
func RegisterInitializer[T Initializer](name string) {}

/*
RegisterValue registers an instance `value` with the `name` as the scope name and returns a function to get parsed `value` from the context. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`.

Example:

	package main

	func main() {
		dbConfig := cfg.RegisterValue[db.Config]("database")

		ctx, err := cfg.Init(context.Background())
		if err != nil {
			panic(err)
		}

		db := database.New(dbConfig(ctx))
	}
*/
func RegisterValue[T any](name string) (getter func(ctx context.Context) *T) {
	return
}

// Init parses configuration from a file, environment or flags. It returns a new context which could be used to retreive values registered with [RegisterValue].
func Init(ctx context.Context) (context.Context, error) {
	return nil, nil
}

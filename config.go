// package cfg implements a framework to load/parse configuration from a file, environment or flags.
package cfg

import "context"

// Initializer is the interface that wraps a instance that provides configuration information and saves the value provided by the configuration.
//
// The `cfg` library parses values from a file, environment or flags and stores into registered instances. If the instance is an Initializer, the method `Initializer.Init()` is called after storing values. The Initializer provider could initialize global instances in this method.
type Initializer interface {
	Init(ctx context.Context) error
}

// RegisterInitializer registers an Initializer instance `value` with the `name` as the scope name. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`. The `value` provider could initialize global instances in `value.Init()`.
func RegisterInitializer(name string, value Initializer) {}

// RegisterValue registers an instance `value` with the `name` as the scope name and returns a function to get parsed `value` from the context. [Init] function parses configuration from a file, environment or flags, stores into `value`, then call `value.Init()`.
func RegisterValue[T any](name string, value *T) (getter func(ctx context.Context) *T) {
	return
}

// Init parses configuration from a file, environment or flags. It returns a new context which could be used to retreive values registered with [RegisterValue].
func Init(ctx context.Context) (context.Context, error) {
	return nil, nil
}

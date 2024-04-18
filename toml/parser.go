package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

type TOML struct{}

func (TOML) ExtNames() []string {
	return []string{"toml"}
}

func (TOML) MimeNames() []string {
	return []string{"application/toml"}
}

func (TOML) TagName() string {
	return "toml"
}

func (TOML) Parse(r io.Reader, a any) error {
	if _, err := toml.NewDecoder(r).Decode(a); err != nil {
		return err
	}
	return nil
}

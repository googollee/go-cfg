package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

type YAML struct{}

func (YAML) ExtNames() []string {
	return []string{"yaml", "yml"}
}

func (YAML) MimeNames() []string {
	return []string{"application/yaml", "application/yml"}
}

func (YAML) TagName() string {
	return "yaml"
}

func (YAML) Parse(r io.Reader, a any) error {
	return yaml.NewDecoder(r).Decode(a)
}

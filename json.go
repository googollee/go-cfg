package cfg

import (
	"encoding/json"
	"io"
)

type JSON struct{}

func (JSON) ExtNames() []string {
	return []string{"json", "js"}
}

func (JSON) MimeNames() []string {
	return []string{"application/javascript", "application/json"}
}

func (JSON) TagName() string {
	return "json"
}

func (JSON) Decode(r io.Reader, a any) error {
	return json.NewDecoder(r).Decode(a)
}

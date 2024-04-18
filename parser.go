package cfg

import (
	"io"
)

type Parser interface {
	ExtNames() []string
	MimeNames() []string
	TagName() string
	Parse(r io.Reader, a any) error
}

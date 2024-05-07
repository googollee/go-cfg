package cfg

import (
	"reflect"
)

type Source interface {
	Setup(t reflect.Type) error
	Parse(v any) error
}

package cfg

import (
	"context"
	"reflect"
)

type Source interface {
	Setup(t reflect.Type) error
	Parse(ctx context.Context, v any) error
}

package cfg

import (
	"context"
	"reflect"
)

type Parser[T any] struct {
	typ     reflect.Type
	sources []Source
}

func Parse[T any](src ...Source) *Parser[T] {
	var v T
	t := reflect.TypeOf(v)

	for _, src := range src {
		if err := src.Setup(t); err != nil {
			panic(err)
		}
	}

	ret := &Parser[T]{
		typ:     t,
		sources: src,
	}
	return ret
}

func (p *Parser[T]) Parse(ctx context.Context, v *T) error {
	for _, src := range p.sources {
		if err := src.Parse(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

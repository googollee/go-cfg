package cfg

import "io"

type Decoder interface {
	Decode(v any) error
}

type NewDecoder struct {
	fn  func(r io.Reader) Decoder
	tag string
}

func F[T Decoder](tag string, fn func(r io.Reader) T) NewDecoder {
	return NewDecoder{
		fn: func(r io.Reader) Decoder {
			return fn(r)
		},
		tag: tag,
	}
}

package cfg

import (
	"encoding"
	"time"
)

type FlagSet interface {
	BoolVar(p *bool, name string, value bool, usage string)
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	Uint64Var(p *uint64, name string, value uint64, usage string)
	UintVar(p *uint, name string, value uint, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	IntVar(p *int, name string, value int, usage string)
	StringVar(p *string, name string, value string, usage string)
	TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string)
	Var(value Value, name string, usage string)
}

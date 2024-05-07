package cfg

import "regexp"

type ConfigSchema interface {
	Schema(schema *Schema)
}

type Schema struct {
	value  any
	fields []Field
}

func (s *Schema) Field(field any) *Field {
	return &Field{}
}

type Field struct {
	name         string
	desc         string
	defaultValue string
	pattern      *regexp.Regexp
}

func (f *Field) Name(name string) *Field               { f.name = name; return f }
func (f *Field) Description(desc string) *Field        { f.desc = desc; return f }
func (f *Field) Default(value string) *Field           { f.defaultValue = value; return f }
func (f *Field) Pattern(pattern *regexp.Regexp) *Field { f.pattern = pattern; return f }

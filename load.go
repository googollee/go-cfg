package cfg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
)

func LoadFile(file string, v any, decoder ...NewDecoder) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("load config(%q) error: %w", file, err)
	}
	defer f.Close()

	return Load(f, v, decoder...)
}

func Load(r io.Reader, v any, decoder ...NewDecoder) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	for _, d := range decoder {
		dt := newStructWithTag(t, d.tag)
		dv := reflect.New(dt)

		if err := d.fn(bytes.NewReader(data)).Decode(dv.Interface()); err == nil {
			copyFields(reflect.ValueOf(v), dv)
			break
		}
	}

	return nil
}

func newStructWithTag(t reflect.Type, tag string) reflect.Type {
	fields := make([]reflect.StructField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		name := field.Tag.Get("cfg")
		field.Tag = reflect.StructTag(fmt.Sprintf("%s:%q", tag, name))
		if field.Type.Kind() == reflect.Struct {
			field.Type = newStructWithTag(field.Type, tag)
		}
		fields = append(fields, field)
	}

	return reflect.StructOf(fields)
}

func copyFields(dst, src reflect.Value) {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}

	for i := 0; i < src.NumField(); i++ {
		if dst.Field(i).Kind() == reflect.Struct {
			copyFields(dst.Field(i), src.Field(i))
			continue
		}
		dst.Field(i).Set(src.Field(i))
	}
}

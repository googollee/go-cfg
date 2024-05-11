package cfg

import (
	"reflect"
	"slices"
)

func walkFields(v reflect.Value, prefix, splitter string, index []int, fn func(key string, index []int, v reflect.Value) error) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		vField := v.Field(i)
		key := prefix + splitter + field.Tag.Get("cfg")
		iField := slices.Clip(append(index, i))

		if field.Type.Kind() == reflect.Struct {
			if err := walkFields(vField, key, splitter, iField, fn); err != nil {
				return err
			}
			continue
		}

		if err := fn(key, iField, vField); err != nil {
			return err
		}
	}

	return nil
}

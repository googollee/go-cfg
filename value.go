package cfg

import (
	"fmt"
	"reflect"
)

func newTypeWithTag(t reflect.Type, fromTag, toTag string) reflect.Type {
	fields := make([]reflect.StructField, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		tagName := field.Tag.Get(fromTag)
		field.Tag = reflect.StructTag(fmt.Sprintf("%s:%q,%s", toTag, tagName, field.Tag))

		if field.Type.Kind() == reflect.Struct {
			field.Type = newTypeWithTag(field.Type, fromTag, toTag)
		}

		fields = append(fields, field)
	}

	return reflect.StructOf(fields)
}

func copyStructField(dst, src reflect.Value) {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	if dst.Kind() == reflect.Ptr {
		dst = dst.Elem()
	}
	t := dst.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		srcField := src.FieldByName(field.Name)
		dstField := dst.FieldByName(field.Name)
		if srcField.Kind() == reflect.Ptr {
			srcField = srcField.Elem()
		}
		if dstField.Kind() == reflect.Ptr {
			dstField = dstField.Elem()
		}

		if field.Type.Kind() == reflect.Struct {
			copyStructField(dstField, srcField)
			continue
		}

		if !srcField.Type().AssignableTo(dstField.Type()) {
			continue
		}

		dstField.Set(srcField)
	}
}

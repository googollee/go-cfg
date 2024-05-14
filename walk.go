package cfg

import (
	"reflect"
	"slices"
	"strings"
)

type fieldMeta struct {
	Key     string
	Default string
	Usage   string
	FullKey string
	Index   []int
}

func tagToFieldMeta(tag string) (meta fieldMeta) {
	sp := strings.SplitN(tag, ",", 3)
	if len(sp) >= 1 {
		meta.Key = strings.TrimSpace(sp[0])
	}
	if len(sp) >= 2 {
		meta.Default = strings.TrimSpace(sp[1])
	}
	if len(sp) >= 3 {
		meta.Usage = strings.TrimSpace(sp[2])
	}

	return
}

func walkFields(v reflect.Value, prefix, splitter string, index []int, fn func(meta fieldMeta, v reflect.Value) error) error {
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
		meta := tagToFieldMeta(field.Tag.Get("cfg"))
		meta.FullKey = prefix + splitter + meta.Key
		meta.Index = slices.Clip(append(index, i))

		if field.Type.Kind() == reflect.Struct {
			if err := walkFields(vField, meta.FullKey, splitter, meta.Index, fn); err != nil {
				return err
			}
			continue
		}

		if err := fn(meta, vField); err != nil {
			return err
		}
	}

	return nil
}

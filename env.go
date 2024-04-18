package cfg

import (
	"io"
	"os"
	"reflect"
	"strings"
)

type Env struct {
	prefix string
}

func (Env) ExtNames() []string {
	return nil
}

func (Env) MimeNames() []string {
	return nil
}

func (Env) TagName() string {
	return ""
}

func (e Env) Parse(r io.Reader, value any) error {
	return parseEnv(value, e.prefix)
}

func parseEnv(value any, prefix string) error {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		tField := t.Field(i)

		if !tField.IsExported() {
			continue
		}

		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			if err := parseEnv(field.Interface(), prefix+"_"+tField.Name); err != nil {
				return err
			}
			return nil
		}

		env := strings.ToUpper(prefix + "_" + tField.Name)
		str, ok := os.LookupEnv(env)
		if !ok {
			continue
		}

		field.SetString(str)
	}

	return nil
}

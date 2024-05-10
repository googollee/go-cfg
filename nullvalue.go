package cfg

import (
	"fmt"
	"reflect"
	"strconv"
)

type flagValue interface {
	Index() []int
	Valid() bool
	CopyTo(value reflect.Value) bool
	UnmarshalText(text []byte) error
}

type nullBase struct {
	index []int
	valid bool
}

func (nb *nullBase) Index() []int { return nb.index }
func (nb *nullBase) Valid() bool  { return nb.valid }

type nullString struct {
	nullBase
	value string
}

func (nv *nullString) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	if value.Kind() != reflect.String {
		return false
	}

	value.Set(reflect.ValueOf(nv.value))
	return true
}

func (nv *nullString) UnmarshalText(text []byte) error {
	nv.valid = true
	nv.value = string(text)
	return nil
}

func (nv *nullString) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}

	return []byte(nv.value), nil
}

type nullInt[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}] struct {
	nullBase
	value T
}

func (nv *nullInt[T]) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	switch value.Kind() {
	case reflect.Int:
		value.Set(reflect.ValueOf(int(nv.value)))
	case reflect.Int8:
		value.Set(reflect.ValueOf(int8(nv.value)))
	case reflect.Int16:
		value.Set(reflect.ValueOf(int16(nv.value)))
	case reflect.Int32:
		value.Set(reflect.ValueOf(int32(nv.value)))
	case reflect.Int64:
		value.Set(reflect.ValueOf(int64(nv.value)))
	}
	return false
}

func (nv *nullInt[T]) Valid() bool {
	return nv.valid
}

func (nv *nullInt[T]) UnmarshalText(text []byte) error {
	nv.valid = true
	return converToInt(string(text), &nv.value)
}

func (nv *nullInt[T]) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%d", nv.value)), nil
}

func converToFloat[T interface {
	~float32 | ~float64
}](str string, v *T) error {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("convert %q to float error: %w", str, err)
	}
	*v = T(f)
	return nil
}

func converToInt[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}](str string, v *T) error {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("convert %q to int error: %w", str, err)
	}
	*v = T(i)
	return nil
}

func converToUint[T interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}](str string, v *T) error {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("convert %q to uint error: %w", str, err)
	}
	*v = T(i)
	return nil
}

func digPtr(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return v.Elem()
}

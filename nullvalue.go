package cfg

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type nullValue interface {
	Init(index []int)
	Index() []int
	Valid() bool
	CopyTo(value reflect.Value) bool
	UnmarshalText(text []byte) error
	MarshalText() ([]byte, error)
}

func newNullValue(kind reflect.Kind) nullValue {
	var ret nullValue

	switch kind {
	case reflect.Int:
		ret = &nullInt[int]{}
	case reflect.Int64:
		ret = &nullInt[int64]{}
	case reflect.String:
		ret = &nullString{}
	}

	return ret
}

type nullBase struct {
	index []int
	valid bool
}

func (nb *nullBase) Index() []int { return nb.index }
func (nb *nullBase) Valid() bool  { return nb.valid }
func (nb *nullBase) Init(index []int) {
	nb.index = index
	nb.valid = false
}

type nullText struct {
	nullBase
	value interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	}
}

func (nv *nullText) CopyTo(value reflect.Value) bool {
	src := reflect.ValueOf(nv.value)
	if !src.Type().AssignableTo(value.Type()) {
		return false
	}

	value.Set(src)
	return true
}

func (nv *nullText) UnmarshalText(text []byte) error {
	nv.valid = true
	return nv.UnmarshalText(text)
}

func (nv *nullText) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}

	return nv.value.MarshalText()
}

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

type nullBool struct {
	nullBase
	value bool
}

func (nv *nullBool) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	if value.Kind() != reflect.Bool {
		return false
	}

	value.Set(reflect.ValueOf(nv.value))
	return true
}

func (nv *nullBool) UnmarshalText(text []byte) error {
	nv.valid = true
	nv.value = false
	switch strings.ToLower(string(text)) {
	case "yes":
		fallthrough
	case "true":
		fallthrough
	case "1":
		nv.value = true
	}
	return nil
}

func (nv *nullBool) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}

	return []byte(fmt.Sprintf("%v", nv.value)), nil
}

type nullDuration struct {
	nullBase
	value time.Duration
}

func (nv *nullDuration) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	if value.Kind() != reflect.Int64 {
		return false
	}

	value.Set(reflect.ValueOf(nv.value))
	return true
}

func (nv *nullDuration) UnmarshalText(text []byte) error {
	nv.valid = true
	var err error
	nv.value, err = time.ParseDuration(string(text))
	return err
}

func (nv *nullDuration) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}

	return []byte(fmt.Sprintf("%s", nv.value)), nil
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

type nullUint[T interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}] struct {
	nullBase
	value T
}

func (nv *nullUint[T]) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	switch value.Kind() {
	case reflect.Int:
		value.Set(reflect.ValueOf(uint(nv.value)))
	case reflect.Int8:
		value.Set(reflect.ValueOf(uint8(nv.value)))
	case reflect.Int16:
		value.Set(reflect.ValueOf(uint16(nv.value)))
	case reflect.Int32:
		value.Set(reflect.ValueOf(uint32(nv.value)))
	case reflect.Int64:
		value.Set(reflect.ValueOf(uint64(nv.value)))
	}
	return false
}

func (nv *nullUint[T]) UnmarshalText(text []byte) error {
	nv.valid = true
	return converToUint(string(text), &nv.value)
}

func (nv *nullUint[T]) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%d", nv.value)), nil
}

type nullFloat[T interface {
	~float32 | ~float64
}] struct {
	nullBase
	value T
}

func (nv *nullFloat[T]) CopyTo(value reflect.Value) bool {
	value = digPtr(value)

	switch value.Kind() {
	case reflect.Float32:
		value.Set(reflect.ValueOf(float32(nv.value)))
	case reflect.Float64:
		value.Set(reflect.ValueOf(float64(nv.value)))
	}
	return false
}

func (nv *nullFloat[T]) UnmarshalText(text []byte) error {
	nv.valid = true
	return converToFloat(string(text), &nv.value)
}

func (nv *nullFloat[T]) MarshalText() ([]byte, error) {
	if !nv.valid {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%g", nv.value)), nil
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

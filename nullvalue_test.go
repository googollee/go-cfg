package cfg

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type testTextValue string

func (v testTextValue) MarshalText() ([]byte, error) {
	return []byte(v), nil
}

func (v *testTextValue) UnmarshalText(text []byte) error {
	*v = testTextValue(string(text))
	return nil
}

type testFlagValue string

func (v testFlagValue) String() string {
	return string(v)
}

func (v *testFlagValue) Set(str string) error {
	*v = testFlagValue(string(str))
	return nil
}

func TestNullValue(t *testing.T) {
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f float64
	var str string
	var dur time.Duration
	var b bool
	var textValue testTextValue
	var flagValue testFlagValue

	tests := []struct {
		value     reflect.Value
		wantType  nullValue
		strValue  string
		wantValue any
	}{
		{
			value:     reflect.ValueOf(&i),
			wantType:  &nullInt[int]{},
			strValue:  "100",
			wantValue: int(100),
		},
		{
			value:     reflect.ValueOf(&i64),
			wantType:  &nullInt[int64]{},
			strValue:  "100",
			wantValue: int64(100),
		},
		{
			value:     reflect.ValueOf(&u),
			wantType:  &nullUint[uint]{},
			strValue:  "100",
			wantValue: uint(100),
		},
		{
			value:     reflect.ValueOf(&u64),
			wantType:  &nullUint[uint64]{},
			strValue:  "100",
			wantValue: uint64(100),
		},
		{
			value:     reflect.ValueOf(&f),
			wantType:  &nullFloat[float64]{},
			strValue:  "100",
			wantValue: float64(100.0),
		},
		{
			value:     reflect.ValueOf(&str),
			wantType:  &nullString{},
			strValue:  "abcd",
			wantValue: "abcd",
		},
		{
			value:     reflect.ValueOf(&b),
			wantType:  &nullBool{},
			strValue:  "true",
			wantValue: true,
		},
		{
			value:     reflect.ValueOf(&dur),
			wantType:  &nullDuration{},
			strValue:  "10s",
			wantValue: 10 * time.Second,
		},
		{
			value:     reflect.ValueOf(&textValue),
			wantType:  &nullText{},
			strValue:  "10s",
			wantValue: testTextValue("10s"),
		},
		{
			value:     reflect.ValueOf(&flagValue),
			wantType:  &nullFlagValue{},
			strValue:  "10s",
			wantValue: testFlagValue("10s"),
		},
	}

	for _, tc := range tests {
		i = 0
		t.Run(fmt.Sprintf("%v", tc.value.Type()), func(t *testing.T) {
			got := newNullValue(tc.value.Type())
			if got, want := reflect.TypeOf(got), reflect.TypeOf(tc.wantType); got != want {
				t.Fatalf("newNullValue(%T) = %v, want: %v", tc.value, got, want)
			}

			if err := got.UnmarshalText([]byte(tc.strValue)); err != nil {
				t.Fatalf("got.UnmarshalText(%q) error: %v", tc.strValue, err)
			}
			strGot, err := got.MarshalText()
			if err != nil {
				t.Fatalf("got.MarshalText() error: %v", err)
			}
			if diff := cmp.Diff(string(strGot), tc.strValue); diff != "" {
				t.Fatalf("diff:\n%s", diff)
			}

			if !got.CopyTo(tc.value) {
				t.Fatalf("got.CopyTo(tc.value) failed")
			}
			if diff := cmp.Diff(tc.value.Elem().Interface(), tc.wantValue); diff != "" {
				t.Fatalf("diff:\n%s", diff)
			}
		})
	}
}

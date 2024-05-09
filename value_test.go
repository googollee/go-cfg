package cfg

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewTypeWithTag(t *testing.T) {
	type Struct struct {
		Str   string `cfg:"str"`
		I     int    `cfg:"i"`
		Inner struct {
			Str string `cfg:"str"`
			I   int    `cfg:"i"`
		} `cfg:"inner"`
	}

	newType := newTypeWithTag(reflect.TypeOf((*Struct)(nil)).Elem(), "cfg", "json")
	jsonData := `{"str":"out_str","i":10,"inner":{"str":"inner_str","i":20}}`
	v := reflect.New(newType).Interface()
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(v); err != nil {
		t.Fatal("decode error:", err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatal("encode error:", err)
	}

	if diff := cmp.Diff(strings.TrimSpace(buf.String()), jsonData); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

func TestCopyStructField(t *testing.T) {
	type CFG struct {
		Str   string `cfg:"str"`
		I     int    `cfg:"i"`
		Inner struct {
			Str string `cfg:"str"`
			I   int    `cfg:"i"`
		} `cfg:"inner"`
	}

	type JSON struct {
		Str   string `json:"str"`
		I     int    `json:"i"`
		Inner struct {
			Str string `json:"str"`
			I   int    `json:"i"`
		} `json:"inner"`
	}

	cfg := CFG{
		Str: "out_str",
		I:   10,
	}
	cfg.Inner.Str = "inner_str"
	cfg.Inner.I = 20

	var j JSON
	copyStructField(reflect.ValueOf(&j), reflect.ValueOf(&cfg))

	want := JSON{
		Str: "out_str",
		I:   10,
	}
	want.Inner.Str = "inner_str"
	want.Inner.I = 20

	if diff := cmp.Diff(j, want); diff != "" {
		t.Errorf("diff:\n%s", diff)
	}
}

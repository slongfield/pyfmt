package pyfmt

import (
	"testing"
)

func TestBasicFormat(t *testing.T) {
	tests := []struct {
		fmtStr string
		params []interface{}
		want   string
	}{
		{"", []interface{}{}, ""},
		{"test", []interface{}{}, "test"},
		{"{{}}", []interface{}{}, "{}"},
		{"{{", []interface{}{}, "{"},
		{"}}", []interface{}{}, "}"},
		{"{}", []interface{}{"test"}, "test"},
		{"{}_{}_{}", []interface{}{"a", "b", "c"}, "a_b_c"},
		{"{1}_{0}", []interface{}{"a", "b"}, "b_a"},
		{"{2}", []interface{}{"a", "b", "c"}, "c"},
		{"{}{1}", []interface{}{"你好", "世界"}, "你好世界"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustFormat(%v, %v) paniced: %v", test.fmtStr, test.params, r)
			}
		}()
		got := MustFormat(test.fmtStr, test.params...)
		if got != test.want {
			t.Errorf("MustFormat(%v, %v) = %v, want %v", test.fmtStr, test.params, got, test.want)
		}
	}
}

func TestBasicMapFormat(t *testing.T) {
	tests := []struct {
		fmtStr string
		params map[string]interface{}
		want   string
	}{
		{"", map[string]interface{}{"test": "test"}, ""},
		{"test", map[string]interface{}{}, "test"},
		{"{{}}", map[string]interface{}{}, "{}"},
		{"{{", map[string]interface{}{}, "{"},
		{"}}", map[string]interface{}{}, "}"},
		{"{test}", map[string]interface{}{"test": "asdf"}, "asdf"},
		{"{a}{c}", map[string]interface{}{"a": "1234", "b": "error", "c": "5678"}, "12345678"},
		{"{hello}{world}", map[string]interface{}{"hello": "你好", "world": "世界"}, "你好世界"},
	}
	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustFormatMap(%v, %v) paniced: %v", test.fmtStr, test.params, r)
			}
		}()
		got := MustFormatMap(test.fmtStr, test.params)
		if got != test.want {
			t.Errorf("MustFormatMap(%v, %v) = %v, want %v", test.fmtStr, test.params, got, test.want)
		}
	}
}

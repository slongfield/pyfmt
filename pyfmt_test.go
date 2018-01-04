package pyfmt

import (
	"testing"
)

func TestBasicFormat(t *testing.T) {
	type ts struct {
		test  string
		hello string
		world string
		a     int8
		b     uint32
		c     int64
	}

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
		{"{}", []interface{}{1}, "1"},
		{"{}", []interface{}{int8(-1)}, "-1"},
		{"{}", []interface{}{uint8(1)}, "1"},
		{"", []interface{}{ts{}}, ""},
		{"{test}", []interface{}{ts{test: "asdf"}}, "asdf"},
		{"{a}{c}", []interface{}{ts{a: 1, b: 2, c: 3}}, "13"},
		{"0b{:b}", []interface{}{3}, "0b11"},
		{"{:#x}", []interface{}{42}, "0x2a"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf(Must("Must({fmtStr}, {params}) paniced: {1}", test, r))
			}
		}()
		got := Must(test.fmtStr, test.params...)
		if got != test.want {
			t.Errorf(Must("Must({fmtStr}, {params}) = {1}, Want: {want}", test, got))
		}
	}
}

// Tests formatting individual values of various types.
// TODO(slongfield): Add more tests.
func TestSingleFormat(t *testing.T) {
	tests := []struct {
		fmtStr string
		param  interface{}
		want   string
	}{
		// String tests
		{"{}", "☺", "☺"},

		// Integer tests
		{"{}", 42, "42"},
		{"{:+#b}", 99, "+0b1100011"},
		{"{: x}", 66, " 42"},

		// Float tests

		// Complex numbers
		{"{}", 0i, "(0+0i)"},
		{"{:3g}", 1 + 1i, "(  1 +1i)"},
		{"{:+12.5g}", 1230000 - 0i, "(   +1.23e+06          +0i)"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf(Must("Must({fmtStr}, {param}) paniced: {1}", test, r))
			}
		}()
		got := Must(test.fmtStr, test.param)
		if got != test.want {
			t.Errorf(Must("Must({fmtStr}, {param}) = {1}, Want: {want}", test, got))
		}
	}

}

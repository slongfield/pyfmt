package pyfmt

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
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
		{"{bar.baz.Bazzle[0]}", []interface{}{pointyMap()}, "1"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Error(Must("Must({fmtStr}, {params}) paniced: {1}", test, r))
			}
		}()
		got := Must(test.fmtStr, test.params...)
		if got != test.want {
			t.Error(Must("Must({fmtStr}, {params}) = {1}, Want: {want}", test, got))
		}
	}
}

type custom int

func (c custom) PyFormat(format string) (string, error) {
	if format == "test" {
		return "test format", nil
	}
	if format == "error" {
		return "", errors.New("Custom formatter error.")
	}
	str := strconv.Itoa(int(c))
	return Fmt("__{}:{}__", format, str)
}

type stringer int

func (s stringer) String() string {
	return "custom stringer"
}

// Tests formatting individual values of various types.
func TestSingleFormat(t *testing.T) {
	tests := []struct {
		fmtStr string
		param  interface{}
		want   string
	}{
		// String tests
		{"{}", "☺", "☺"},
		{"{:}", "0", "0"},
		{"{:t}", "", "string"},
		{"asdf{:10}", "1234", "asdf      1234"},
		{"{:💩^10}", "poop", "💩💩💩poop💩💩💩"},

		// Integer tests
		{"{}", 42, "42"},
		{"{:+#b}", 99, "+0b1100011"},
		{"{: x}", 66, " 42"},
		{"{:t}", 66, "int"},
		{"{:^10}", 1, "    1     "},
		{"{:^10}", 10, "    10    "},
		{"{:^10}", 100, "   100    "},
		{"{:^10}", 1000, "   1000   "},
		{"{:<10}", 1, "1         "},
		{"{:<10}", 10, "10        "},
		{"{:<10}", 100, "100       "},
		{"{:<10}", 1000, "1000      "},
		{"{:<10}", 1, "1         "},
		{"{:>10}", 1, "         1"},
		{"{:>10}", 10, "        10"},
		{"{:>10}", 100, "       100"},
		{"{:>10}", 1000, "      1000"},
		{"{:<10x}", -10, "-a        "},
		{"{:+d}", 99, "+99"},
		{"{:=4d}", -99, "- 99"},
		{"{:=4d}", 99, "  99"},
		{"{:=+4d}", 99, "+ 99"},
		{"{:#010X}", 0, "0X00000000"},
		{"{:#9b}", -3, "    -0b11"},
		{"0{:#0b}", 0, "00b0"},
		{"{:#0d}", 0, "0"},
		{"{::<010X}", 0, "0:::::::::"},
		{"{: #b}", 0, " 0b0"},
		{"{: #0b}", 0, " 0b0"},
		{"{: #0b}", -1, "-0b1"},
		{"{: 010X}", 0, " 000000000"},
		{"{: 010X}", -10, "-00000000A"},
		{"{::>+9X}", 1234, ":::::+4D2"},
		{"{::=#10X}", -1, "-0X::::::1"},
		{"{:10X}", 0, "         0"},

		// Float tests
		{"{:.0%}", 0.25, "25%"},
		{"{:g}", math.Inf(+1), "+Inf"},
		{"{:g}", math.Inf(-1), "-Inf"},
		// No negative zero in Go constants
		{"{:g}", math.Copysign(-0.0, -1), "-0"},
		{"{:g}", math.NaN(), "NaN"},
		{"{:.1%}", 0.25, "25.0%"},
		{"{:.3%}", 0.0, "0.000%"},
		{"{:.0%}", -2.0, "-200%"},
		{"{:.3%}", 1.2, "120.000%"},
		{"{::<20E}", 0.0, "0.000000E+00::::::::"},
		{"{:<1.0%}", math.Copysign(-0.0, -1), "-0%"},
		{"{:<1.0%}", -0.1, "-10%"},
		{"{::<#1.0E}", 0.0, "0E+00"},
		{"{: 8.1E}", 1.1, " 1.1E+00"},
		{"{: 01.1E}", 1.9, " 1.9E+00"},

		// Complex numbers
		{"{}", 0i, "(0+0i)"},
		{"{:3g}", 1 + 1i, "(  1 +1i)"},
		{"{:+12.5g}", 1230000 - 0i, "(   +1.23e+06          +0i)"},

		// Structs
		{"{}", struct {
			a int
			b int
		}{1, 2}, "{1 2}"},
		{"{:r}", struct {
			a int
			b int
		}{1, 2}, "struct { a int; b int }{a:1, b:2}"},
		{"{:s}", struct {
			a int
			b int
		}{1, 2}, "{a:1 b:2}"},

		// Custom formatter
		{"{:test}", custom(99), "test format"},
		{"{:asdf}", custom(3), "__asdf:3__"},
		{"{0[0]:1234}", []custom{99}, "__1234:99__"},
		{"{Test:test}", struct{ Test custom }{Test: 3}, "test format"},
		{"{Test:1234}", struct{ Test custom }{Test: 3}, "__1234:3__"},
		// Custom formatters don't work in unexported struct variables
		{"{test}", struct{ test custom }{test: 1234}, "1234"},

		// Custom fmt.Stringer
		{"{}", stringer(3), "custom stringer"},
		{"{0[0]}", []stringer{99}, "custom stringer"},
		{"{Test}", struct{ Test stringer }{Test: 42}, "custom stringer"},
		// Custom stringers don't work in unexported struct variables.
		{"{test}", struct{ test stringer }{test: 6789}, "6789"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Error(Must("Must({fmtStr}, {param}) paniced: {1}", test, r))
			}
		}()
		got := Must(test.fmtStr, test.param)
		if got != test.want {
			t.Error(Must("Must({fmtStr}, {param}) = {1}, Want: {want}", test, got))
		}
	}
}

func BenchmarkPrintEmptyParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Must("")
		}
	})
}

func BenchmarkFmtForComparison(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = fmt.Sprintf("")
		}
	})
}

func BenchmarkCenteredParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Must("{:^100}", "test")
		}
	})
}

func BenchmarkLargeString(b *testing.B) {
	test := strings.Repeat("{0}", 1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Must(test, "test")
		}
	})
}

func BenchmarkFmtLargeString(b *testing.B) {
	test := strings.Repeat("%[0]v", 1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = fmt.Sprintf(test, "test")
		}
	})
}

func BenchmarkComplexFormat(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Must("{0[0]:😄^+#30.30b}", []int{42})
		}
	})
}

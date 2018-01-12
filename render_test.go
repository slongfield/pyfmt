package pyfmt

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestSplitFlags(t *testing.T) {
	tests := []string{"", "4<", "+=", "^10.3", ":> #010.4X", "<0%", "10.10E", "#x"}
	var flagPattern = regexp.MustCompile(`\A((?:.[<>=^])|(?:[<>=^])?)([\+\- ]?)(#?)(0?)(\d*)(\.\d*)?([bdoxXeEfFgGrts%]?)\z`)

	for _, test := range tests {
		align, sign, radix, zeroPad, minWidth, precision, verb, err := splitFlags(test)

		if err != nil {
			t.Error(Error("splitFlags({}) errored: {}!", test, err))
		}

		if !flagPattern.MatchString(test) {
			t.Error(Error("Could not match with regex!: {}", test))
		}

		got := []string{test, align, sign, radix, zeroPad, minWidth, precision, verb}
		want := flagPattern.FindStringSubmatch(test)
		if !reflect.DeepEqual(got, want) {
			t.Error(Error("splitFlags({}) = {} Want: {}", test, got, want))
		}
	}
}

func TestSplitFlagsError(t *testing.T) {
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		flagStr string
		want    flags
	}{
		{"", flags{renderVerb: "v"}},
		{">>", flags{fillChar: '>', align: right, renderVerb: "v"}},
		{">10.10", flags{align: right, minWidth: "10", precision: ".10", renderVerb: "v"}},
		{"#x", flags{showRadix: true, renderVerb: "x"}},
		{"#X", flags{showRadix: true, renderVerb: "X"}},
		// Neg sign doesn't get picked up.
		{"-.4o", flags{precision: ".4", sign: "", renderVerb: "o"}},
		{"+.4o", flags{precision: ".4", sign: "+", renderVerb: "o"}},
		{"r", flags{renderVerb: "#v"}},
		{"#010X", flags{showRadix: true, align: padSign, fillChar: '0', minWidth: "10", renderVerb: "X"}},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Error(Error("parseFlags({flagStr}) paniced: {1}", test, r))
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err != nil {
			t.Error(Error("parseFlags({flagStr}) Errored: {1}", test, err))
		}
		if !reflect.DeepEqual(test.want, r.flags) {
			t.Error(Error("parseFlags({flagStr}) Got: \n{1:s} Want \n{want:s}", test, r.flags))
		}
	}
}

func TestParseFlagsError(t *testing.T) {
	tests := []struct {
		flagStr string
		want    string
	}{
		{"asdf", "Invalid"},
		{":::", "Invalid"},
		{">10.10.", "Invalid"},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Error(Error("parseFlags({flagStr}) paniced: {1}", test, r))
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err == nil {
			t.Error(Error("parseFlags({flagStr}) did not raise an error", test))
		}
		if !strings.Contains(err.Error(), test.want) {
			t.Error(Error("parseFlags({flagStr}) raised {1}, missing want string {want}", test, err))
		}
	}
}

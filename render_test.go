package pyfmt

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

const flagRegex = `\A((?:.[<>=^])|(?:[<>=^])?)([\+\- ]?)(#?)(0?)(\d*)(\.\d*)?([bdoxXeEfFgGrts%]?)\z`

func TestSplitFlags(t *testing.T) {
	var flagPattern = regexp.MustCompile(flagRegex)

	tests := []string{"", "4<", "+=", "^10.3", ":> #010.4X",
		"<0%", "10.10E", "#x", "<<", "==", "💩<"}

	for _, test := range tests {
		align, sign, radix, zeroPad, minWidth, precision, verb, err := splitFlags(test)

		if err != nil {
			t.Error(Must("splitFlags({}) errored: {}!", test, err))
		}

		if !flagPattern.MatchString(test) {
			t.Error(Must("Could not match with regex!: {}", test))
		}

		got := []string{test, align, sign, radix, zeroPad, minWidth, precision, verb}
		want := flagPattern.FindStringSubmatch(test)
		if !reflect.DeepEqual(got, want) {
			t.Error(Must("splitFlags({}) = \n{:r} Want: \n{:r}", test, got, want))
		}
	}
}

func TestSplitFlagsError(t *testing.T) {
	tests := []string{"<><>", "asdf", "^^^", "^#xx", ":>  #010.4x"}
	for _, test := range tests {
		_, _, _, _, _, _, _, err := splitFlags(test)
		if err == nil {
			t.Error(Must("splitFlags({}) did not error!", test))
		}
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		flagStr string
		want    flags
	}{
		{"", flags{renderVerb: "v", empty: true}},
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
				t.Error(Must("parseFlags({flagStr}) paniced: {1}", test, r))
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err != nil {
			t.Error(Must("parseFlags({flagStr}) Errored: {1}", test, err))
		}
		if !reflect.DeepEqual(test.want, r.flags) {
			t.Error(Must("parseFlags({flagStr}) Got: \n{1:s} Want \n{want:s}", test, r.flags))
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
				t.Error(Must("parseFlags({flagStr}) paniced: {1}", test, r))
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err == nil {
			t.Error(Must("parseFlags({flagStr}) did not raise an error", test))
		}
		if !strings.Contains(err.Error(), test.want) {
			t.Error(Must("parseFlags({flagStr}) raised {1}, missing want string {want}", test, err))
		}
	}
}

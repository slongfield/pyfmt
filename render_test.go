package pyfmt

import (
	"reflect"
	"strings"
	"testing"
)

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

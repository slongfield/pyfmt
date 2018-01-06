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
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Error("parseFlags(%v) paniced: %v", test.flagStr, r)
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err != nil {
			t.Error("parseFlags(%v) Errored: %v", test.flagStr, err)
		}
		if !reflect.DeepEqual(test.want, r.flags) {
			t.Error("parseFlags(%v) Got: %v Want %v", test.flagStr, r.flags, test.want)
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
				t.Error("parseFlags(%v) paniced: %v", test.flagStr, r)
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err == nil {
			t.Error("parseFlags(%v) did not raise an error", test.flagStr)
		}
		if !strings.Contains(err.Error(), test.want) {
			t.Error("parseFlags(%v) raised %v, missing want string %v", test.flagStr, err, test.want)
		}
	}

}

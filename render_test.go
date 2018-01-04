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
		{"", flags{}},
		{">>", flags{fillChar: '>', align: right}},
		{">10.10", flags{align: right, minWidth: "10", precision: ".10"}},
		{"#x", flags{showRadix: true, renderType: hex}},
		{"#X", flags{showRadix: true, renderType: hexCap}},
		// Neg sign doesn't get picked up.
		{"-.4o", flags{precision: ".4", sign: "", renderType: octal}},
		{"+.4o", flags{precision: ".4", sign: "+", renderType: octal}},
	}

	for _, test := range tests {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("parseFlags(%v) paniced: %v", test.flagStr, r)
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err != nil {
			t.Errorf("parseFlags(%v) Errored: %v", test.flagStr, err)
		}
		if !reflect.DeepEqual(test.want, r.flags) {
			t.Errorf("parseFlags(%v) Got: %v Want %v", test.flagStr, r.flags, test.want)
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
				t.Errorf("parseFlags(%v) paniced: %v", test.flagStr, r)
			}
		}()
		r := render{}
		err := r.parseFlags(test.flagStr)
		if err == nil {
			t.Errorf("parseFlags(%v) did not raise an error", test.flagStr)
		}
		if !strings.Contains(err.Error(), test.want) {
			t.Errorf("parseFlags(%v) raised %v, missing want string %v", test.flagStr, err, test.want)
		}
	}

}

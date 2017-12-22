package pyfmt

import (
	"reflect"
	"testing"
)

func TestGetElement(t *testing.T) {
	tests := []struct {
		elems        []interface{}
		lookupStr    string
		lookupOffset int
		want         interface{}
	}{
		{[]interface{}{3}, "", 0, 3},
		{[]interface{}{3, "asdf"}, "", 1, "asdf"},
		{[]interface{}{3}, "0", 0, 3},
		{[]interface{}{3, "asdf"}, "1", 0, "asdf"},
		{[]interface{}{struct{ Test string }{Test: "asdf"}}, "Test", 0, "asdf"},
	}

	for _, test := range tests {
		got, err := getElement(test.lookupStr, test.lookupOffset, test.elems...)
		if err != nil {
			t.Errorf("getElement(%v, %v, %v) Errored: %v", test.lookupStr, test.lookupOffset, test.elems, err)
		}
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("getElement(%v, %v, %v) = %v Want: %v", test.lookupStr, test.lookupOffset, test.elems, got, test.want)
		}
	}
}

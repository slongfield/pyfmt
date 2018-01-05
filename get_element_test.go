package pyfmt

import (
	"reflect"
	"testing"
)

type inner struct {
	test int64
}

type outer struct {
	first  inner
	second inner
}

func nestedMap() map[string]map[string]map[string]int64 {
	m := make(map[string]map[string]map[string]int64)
	m["test"] = make(map[string]map[string]int64)
	m["test"]["bar"] = make(map[string]int64)
	m["test"]["bar"]["foo"] = 99
	return m
}

// elementFromValue will try to turn a reflect.Value into an interface{} when possible.
func elementFromValue(val reflect.Value) interface{} {
	if !val.IsValid() {
		return nil
	}
	if val.CanInterface() {
		return val.Interface()
	}
	// TODO(slongfield): Get a larger set of values.
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uint64(val.Uint())
	case reflect.String:
		return val.String()
	}
	return val
}

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
		{[]interface{}{struct{ test string }{test: "asdf"}}, "test", 0, "asdf"},
		{[]interface{}{map[string]int{"test": 3, "asdf": 4}}, "test", 0, 3},
		{[]interface{}{map[string]int{"test": 3, "asdf": 4}}, "asdf", 0, 4},
		{[]interface{}{3, []string{"foo", "bar"}}, "1[1]", 0, "bar"},
		{[]interface{}{struct{ foo outer }{foo: outer{first: inner{test: 5}}}}, "foo.first.test", 0, int64(5)},
		{[]interface{}{nestedMap()}, "test[bar].foo", 0, int64(99)},
	}

	for _, test := range tests {
		got, err := getElement(test.lookupStr, test.lookupOffset, test.elems...)
		if err != nil {
			t.Errorf("getElement(%v, %v, %v) Errored: %v", test.lookupStr, test.lookupOffset, test.elems, err)
		}
		// If we got a reflect.Value, pull out the underlying element. These print correctly, but
		// reflect.DeepEqual doesn't like the unboxing.
		switch got.(type) {
		case reflect.Value:
			got = elementFromValue(got.(reflect.Value))
		}
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("getElement(%v, %v, %v) = %v Want: %v", test.lookupStr, test.lookupOffset, test.elems, got, test.want)
		}
	}
}

func TestGetElementErrors(t *testing.T) {
	tests := []struct {
		elems        []interface{}
		lookupStr    string
		lookupOffset int
	}{
		{[]interface{}{}, "", 1},
		{[]interface{}{}, "test", 0},
		{[]interface{}{3}, "", 1},
		{[]interface{}{3}, "1", 0},
		{[]interface{}{3, "asdf"}, "5", 0},
		{[]interface{}{struct{ Test string }{Test: "asdf"}}, "Best", 0},
		{[]interface{}{struct{ test string }{test: "asdf"}}, "asdf", 0},
		{[]interface{}{map[string]int{"test": 3, "asdf": 4}}, "jkl;", 0},
		{[]interface{}{map[int]string{3: "test", 4: "asdf"}}, "3", 0},
		{[]interface{}{3, []string{"foo", "bar"}}, "1[5]", 0},
	}

	for _, test := range tests {
		_, err := getElement(test.lookupStr, test.lookupOffset, test.elems...)
		if err == nil {
			t.Errorf(Must("getElement({lookupStr}, {lookupOffset}, {elems}) Did not error!", test))
		}
	}
}

func TestSplitName(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{"", []string{}},
		{"test", []string{"test"}},
		{"foo.bar", []string{"foo", "bar"}},
		{"foo[3].bar", []string{"foo", "3", "bar"}},
		{"baz[3][4][5]", []string{"baz", "3", "4", "5"}},
		{"bar[3].4[5]", []string{"bar", "3", "4", "5"}},
	}

	for _, test := range tests {
		got, err := splitName(test.name)
		if err != nil {
			t.Errorf(Must("splitName({name}) Errored: {1}", test, err))
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(Must("splitName({name}) = {1} Want: {want}", test, got))
		}
	}
}

func TestSplitNameError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"te[3]st"},
		{"["},
		{"[[["},
		{"]"},
		{"]]]"},
	}

	for _, test := range tests {
		_, err := splitName(test.name)
		if err == nil {
			t.Errorf(Must("splitName({name}) did not error", test))
		}
	}
}

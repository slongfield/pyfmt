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

type outptr struct {
	ptr *inner
}

type b struct {
	Bazzer int
	Bazzle []int
}

func nestedMap() map[string]map[string]map[string]int64 {
	m := make(map[string]map[string]map[string]int64)
	m["test"] = make(map[string]map[string]int64)
	m["test"]["bar"] = make(map[string]int64)
	m["test"]["bar"]["foo"] = 99
	return m
}

// pointyMap is a map with some pointers and lists in it.
func pointyMap() interface{} {
	return map[string]interface{}{
		"bar": map[string]interface{}{
			"baz":  &b{0, []int{1, 2, 3}},
			"buzz": []int{4, 5, 6}},
		"baz":    []int{7, 8, 9},
		"bazzle": []string{"10", "11", "12"}}
}

// elementFromValue will try to turn a reflect.Value into an interface{} when possible.
func elementFromValue(val reflect.Value) interface{} {
	if !val.IsValid() {
		return nil
	}
	if val.CanInterface() {
		return val.Interface()
	}
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
		{[]interface{}{outptr{ptr: &inner{test: 3}}}, "ptr.test", 0, int64(3)},
		{[]interface{}{pointyMap()}, "bar.baz.Bazzle[0]", 0, int(1)},
		{[]interface{}{pointyMap()}, "bazzle", 0, []string{"10", "11", "12"}},
		{[]interface{}{pointyMap()}, "bazzle[1]", 0, "11"},
	}

	for _, test := range tests {
		got, err := getElement(test.lookupStr, test.lookupOffset, test.elems...)
		if err != nil {
			t.Error(Must("getElement({lookupStr}, {lookupOffset}, {elems}) Errored: {1}", test, err))
		}
		// If we got a reflect.Value, pull out the underlying element. These print correctly, but
		// reflect.DeepEqual doesn't like the unboxing.
		switch got.(type) {
		case reflect.Value:
			got = elementFromValue(got.(reflect.Value))
		}
		if !reflect.DeepEqual(test.want, got) {
			t.Error(Must("getElement({lookupStr}, {lookupOffset}, {elems}) = {1} ({1:t}) Want: {want} ({want:t})", test, got))
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
			t.Error(Must("getElement({lookupStr}, {lookupOffset}, {elems}) Did not error!", test))
		}
	}
}

func TestSplitName(t *testing.T) {
	tests := []struct {
		name      string
		wantfield string
		wantrem   string
	}{
		{"", "", ""},
		{"test", "test", ""},
		{"foo.bar", "foo", "bar"},
		{"foo[3].bar", "foo", "[3].bar"},
		{"[3].bar", "3", "bar"},
		{"bar", "bar", ""},
		{"baz[3][4][5]", "baz", "[3][4][5]"},
		{"bar[3].4[5]", "bar", "[3].4[5]"},
		{"[3].4[5]", "3", "4[5]"},
		{"4[5]", "4", "[5]"},
		{"[5]", "5", ""},
	}

	for _, test := range tests {
		got, rem, err := splitName(test.name, false)
		if err != nil {
			t.Error(Must("splitName({name}) Errored: {1}", test, err))
		}
		if !reflect.DeepEqual(got, test.wantfield) {
			t.Error(Must("splitName({name}) = {1}, {2} Want: {wantfield}, {wantrem}", test, got, rem))
		}
		if !reflect.DeepEqual(rem, test.wantrem) {
			t.Error(Must("splitName({name}) = {1}, {2} Want: {wantfield}, {wantrem}", test, got, rem))
		}
	}
}

func TestSplitNameError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"[3]test"},
		{"["},
		{"[[["},
		{"]"},
		{"]]]"},
	}

	for _, test := range tests {
		_, _, err := splitName(test.name, false)
		if err == nil {
			t.Error(Must("splitName({name}) did not error", test))
		}
	}
}

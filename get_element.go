package pyfmt

import (
	"fmt"
	"reflect"
	"strconv"
)

// getElement takes a string, an offset, and a list of elements, and returns the element in the
// list that matches the request. The rules for matching are:
// - if 'name' does not contain a '.' or '[]':
//   - if 'name' is "", use the offset, and return the element at that offset.
//   - if 'name' is a digit, use that as an offset
//   - if 'name' is a string, 'elems' is exactly one element long, and the first element of elems
//      is a map, use the string as a lookup into the map. If the first element is a struct, get
//      the member of the struct with that name.
// - if the 'name' contains '[{id}]', and the 'id' is a text, use the part of the string before
//   '[{id}]' to look up the name as above,
// - if the 'name' contains a '.' use the above lookup rules on the part before the dot to look up
//   an element, and then follow the rules as above.
//
func getElement(name string, offset int, elems ...interface{}) (interface{}, error) {
	if len(elems) == 0 {
		return nil, fmt.Errorf("attempted to fetch %v/%v from empty list", name, offset)
	}
	if name == "" {
		if offset < len(elems) {
			return elems[offset], nil
		}
		return nil, fmt.Errorf("too large offset: %v", offset)
	}
	if parse, err := strconv.ParseUint(name, 10, 64); err == nil {
		if parse < uint64(len(elems)) {
			return elems[parse], nil
		}
		return nil, fmt.Errorf("too large parse: %v", parse)
	}

	return getElementByName(name, elems[0])
}

// Gets the element by name, if possible.
func getElementByName(name string, src interface{}) (interface{}, error) {
	if reflect.ValueOf(src).Kind() == reflect.Struct {
		v := reflect.ValueOf(src).FieldByName(name)
		if v.IsValid() {
			return getElementFromValue(v), nil
		}
		return nil, fmt.Errorf("Could not find field: %s", name)
	}
	return nil, nil
}

func getElementFromValue(val reflect.Value) interface{} {
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
	return nil
}

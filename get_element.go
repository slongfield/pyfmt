package pyfmt

import "fmt"

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
	// TODO(slongfield): Fill in and test, then use to replace all the Format methods.
	if name == "" {
		if offset < len(elems) {
			return elems[offset], nil
		}
		return nil, fmt.Errorf("too large offset: %v", offset)
	}
	return nil, nil
}

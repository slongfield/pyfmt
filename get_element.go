package pyfmt

import (
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
func getElement(name string, offset int, elems ...interface{}) (interface{}, error) {
	if len(elems) == 0 {
		return nil, Error("attempted to fetch {}/{} from empty list", name, offset)
	}
	field, remainder, err := splitName(name, true)
	if err != nil {
		return nil, err
	}

	val := elems[0]
	found := false
	if field == "" {
		if offset < len(elems) {
			val = elems[offset]
			found = true
		} else {
			return nil, Error("too large offset: {}", offset)
		}
	} else if parse, err := strconv.ParseUint(field, 10, 64); err == nil {
		if parse < uint64(len(elems)) {
			val = elems[parse]
			found = true
		} else {
			return nil, Error("index out of bounds: {}", parse)
		}
	}

	if !found {
		val, err = elementByName(field, val)
		if err != nil {
			return nil, err
		}
	}
	for remainder != "" {
		field, remainder, err = splitName(remainder, false)
		if err != nil {
			return nil, err
		}
		val, err = elementByName(field, val)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

// splitName splits the first subfield off of the name, returning both the that was split off and
// the remainder. Errors if it can't be split.  Note that this does treat test[foo].bar and
// test.bar[foo] as being interchangeable. This normally makes sense, especially for
// structs-of-structs and maps-of-maps, but may be somewhat strange for lists of lists, where
// a[5][6] can be written a.5.6.
func splitName(name string, first bool) (string, string, error) {
	end := len(name)
	foundOpen := false
	for i := 0; i < end; {
		cachei := i
		for i < end && !(name[i] == '[' || name[i] == '.' || name[i] == ']') {
			i++
		}
		if i < end {
			if name[i] == '.' {
				return name[:i], name[(i + 1):], nil
			}
			if name[i] == ']' && foundOpen {
				if i+1 < end && !(name[i+1] == '[' || name[i+1] == '.') {
					return "", "", Error("must begin a new subfield after a closing bracket in {}", name)
				}
				if i+1 < end {
					return name[cachei:i], name[(i + 2):], nil
				}
				return name[cachei:i], name[(i + 1):], nil
			}
			if name[i] == ']' && !foundOpen {
				return "", "", Error("unmatched ] in {}", name)
			}
			if name[i] == '[' && foundOpen {
				return "", "", Error("unmatched [ in {}", name)
			}
			if name[i] == '[' && i == 0 && !first {
				foundOpen = true
				i++
				cachei = i
				continue
			}
			if name[i] == '[' {
				return name[:i], name[i:], nil
			}
		} else {
			return name, "", nil
		}
		i++
	}
	if foundOpen {
		return "", "", Error("unmatched [ in {}", name)
	}
	return name, "", nil
}

// elementByName will get the element by name if it's a struct or map, the an element by number
// from an Array or Slice, and error out otherwise. If possible, will return an interface{} value,
// but may return a reflect.Value if it cannot be interfaced (e.g., for unexported struct fields)
func elementByName(name string, src interface{}) (interface{}, error) {
	var srcVal reflect.Value
	switch src.(type) {
	case reflect.Value:
		srcVal = src.(reflect.Value)
	default:
		srcVal = reflect.ValueOf(src)
	}

	switch srcVal.Kind() {
	case reflect.Ptr:
		if srcVal.IsNil() {
			return nil, Error("attempted to dereference nil pointer {}", name)
		}
		return elementByName(name, reflect.Indirect(srcVal))
	case reflect.Struct:
		v := srcVal.FieldByName(name)
		if v.IsValid() {
			if v.CanInterface() {
				return v.Interface(), nil
			}
			return v, nil
		}
		return nil, Error("could not find field: {}", name)
	case reflect.Map:
		if !(reflect.ValueOf(name).Type().AssignableTo(srcVal.Type().Key())) {
			return nil, Error("could not look up key {} from map {}", name, src)
		}
		v := srcVal.MapIndex(reflect.ValueOf(name))
		if v.IsValid() {
			if v.CanInterface() {
				return v.Interface(), nil
			}
			return v, nil
		}
		return nil, Error("could not find key: {}", name)
	case reflect.Array, reflect.Slice:
		if parse, err := strconv.ParseUint(name, 10, 64); err == nil {
			if parse < uint64(srcVal.Len()) {
				v := srcVal.Index(int(parse))
				if v.IsValid() {
					if v.CanInterface() {
						return v.Interface(), nil
					}
					return v, nil
				}
				return nil, Error("could not get index: {}", name)
			}
			return nil, Error("index out of bounds: {}", parse)
		}
		return nil, Error("could not parse index: {}", name)
	default:
		return nil, Error("attempted to get item by name from non-struct, non-map: {} {}", src, srcVal.Kind())
	}
}

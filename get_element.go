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
	fields, err := splitName(name)
	if err != nil {
		return nil, err
	}

	val := elems[0]
	idx := 0
	if len(fields) == 0 || fields[0] == "" {
		if offset < len(elems) {
			val = elems[offset]
			idx++
		} else {
			return nil, Error("too large offset: {}", offset)
		}
	} else if parse, err := strconv.ParseUint(fields[0], 10, 64); err == nil {
		if parse < uint64(len(elems)) {
			val = elems[parse]
			idx++
		} else {
			return nil, Error("index out of bounds: {}", parse)
		}
	}

	if idx < len(fields) {
		for _, field := range fields[idx:] {
			val, err = elementByName(field, val)
			if err != nil {
				return nil, err
			}
		}
	}
	return val, nil
}

// splitName splits out the name into the subfields. Errors if it can't cleanly split.
// Note that this does treat test[foo].bar and test.bar[foo] as being interchangeable. This normally
// makes sense, especially for structs-of-structs and maps-of-maps, but may be somewhat strange for
// lists of lists, where a[5][6] can be written a.5.6.
func splitName(name string) ([]string, error) {
	subNames := make([]string, 0, 0)
	end := len(name)
	foundOpen := false
	for i := 0; i < end; {
		cachei := i
		for i < end && !(name[i] == '[' || name[i] == '.' || name[i] == ']') {
			i++
		}
		if i < end {
			if name[i] == '.' {
				subNames = append(subNames, name[cachei:i])
				i++
				cachei = i
				continue
			}
			if name[i] == ']' && foundOpen {
				subNames = append(subNames, name[cachei:i])
				i++
				cachei = i
				foundOpen = false
				if i < end && !(name[i] == '[' || name[i] == '.') {
					return nil, Error("must begin a new subfield after a closing bracket in {}", name)
				} else if i < end {
					foundOpen = (name[i] == '[')
					i++
				}
				continue
			}
			if name[i] == ']' && !foundOpen {
				return nil, Error("unmatched ] in {}", name)
			}
			if name[i] == '[' {
				subNames = append(subNames, name[cachei:i])
				foundOpen = true
				i++
				cachei = i
				continue
			}
		} else {
			subNames = append(subNames, name[cachei:i])
		}
		i++
	}
	if foundOpen {
		return nil, Error("unmatched [ in {}", name)
	}
	return subNames, nil
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

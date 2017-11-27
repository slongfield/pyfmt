package pyfmt

import (
	"errors"
	"fmt"
)

// Using a simple []byte instead of bytes.Buffer to avoid the dependency.
type buffer []byte

func (b *buffer) WriteString(s string) {
	*b = append(*b, s...)
}

// ff is used to store a formatter's state.
type ff struct {
	buf buffer

	// argList is the list of arguments, if it was passed that way.
	argList []interface{}
	useList bool
	listPos int

	// argMap is a map of strings, as an alternate format parameter method
	argMap map[string]interface{}
}

// newFormater creates a new ff struct.
// TODO(slongfield): Investigate using a sync.Pool to avoid reallocation.
func newFormater() *ff {
	return &ff{}
}

// doFormat parses the string, and executes a format command. Stores the output in ff's buf.
func (f *ff) doFormat(format string) error {
	end := len(format)
	for i := 0; i < end; {
		cachei := i
		// First, get to a '{'
		for i < end && format[i] != '{' {
			// If we see a '}' before a '{' it's an error, unless the next character is also a '}'.
			if format[i] == '}' {
				if i+1 == end || format[i+1] != '}' {
					return errors.New("Single '}' encountered in format string")
				} else {
					f.buf.WriteString(format[cachei:i])
					i++
					cachei = i
				}
			}
			i++
		}
		if i > cachei {
			f.buf.WriteString(format[cachei:i])
		}
		if i >= end {
			break
		}
		i++
		// If the next character is also '{', just put the '{' back in and continue.
		if format[i] == '{' {
			f.buf.WriteString("{")
			i++
		}
		// TODO(slongfield): Inner loop to parse the format string.
	}
	return nil
}

// Format is the equivlent of Python's string.format() function. Takes a list of possible elements
// to use in formatting, and substitutes them. Only allows for the {}, {0} style of substitutions.
func Format(format string, a ...interface{}) (string, error) {
	f := newFormater()
	f.argList = a
	f.useList = true
	err := f.doFormat(format)
	if err != nil {
		return "", err
	}
	s := string(f.buf)
	return s, nil
}

// FormatMap is similar to Python's string.format(), but takes a map from name to interface to allow
// for {name} style formatting.
func FormatMap(format string, a map[string]interface{}) (string, error) {
	return fmt.Sprintf(format, a), nil
}

// MustFormat is like Format, but panics on error.
func MustFormat(format string, a ...interface{}) string {
	s, err := Format(format, a...)
	if err != nil {
		panic(err)
	}
	return s
}

// MustFormatMap is like FormatMap, but panics on error.
func MustFormatMap(format string, a map[string]interface{}) string {
	s, err := FormatMap(format, a)
	if err != nil {
		panic(err)
	}
	return s
}

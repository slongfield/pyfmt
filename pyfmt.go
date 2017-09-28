package pyfmt

import "fmt"

// Format is the equivlent of Python's string.format() function. Takes a list of possible elements
// to use in formatting, and substitutes them. Only allows for the {}, {0} style of substitutions.
func Format(format string, a ...interface{}) (string, error) {
	return fmt.Sprintf(format, a...), nil
}

// FormatMap is similar to Python's string.format(), but takes a map from name to interface to allow
// for {name} style formatting.
func FormatMap(format string, a map[string]interface{}) (string, error) {
	return fmt.Sprintf(format, a), nil
}

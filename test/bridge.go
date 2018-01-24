package main

import "C"

import (
	"fmt"

	"github.com/slongfield/pyfmt"
)

// FormatOneInt takes a format string and a single int, and formats it.
//export FormatOneInt
func FormatOneInt(cformat *C.char, arg C.int) *C.char {
	result, err := pyfmt.Fmt(C.GoString(cformat), int32(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v\n", err)
	}
	return C.CString(result)
}

// FormatOneFloat takes a format string and a single float, and formats it.
//export FormatOneFloat
func FormatOneFloat(cformat *C.char, arg C.float) *C.char {
	result, err := pyfmt.Fmt(C.GoString(cformat), float32(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v\n", err)
	}
	return C.CString(result)
}

// FormatOneDouble takes a format string and a single double, and formats it.
//export FormatOneDouble
func FormatOneDouble(cformat *C.char, arg C.double) *C.char {
	result, err := pyfmt.Fmt(C.GoString(cformat), float64(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v\n", err)
	}
	return C.CString(result)
}

// FormatOneString takes a format string and a single string, and formats it.
//export FormatOneString
func FormatOneString(cformat *C.char, arg *C.char) *C.char {
	result, err := pyfmt.Fmt(C.GoString(cformat), C.GoString(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v with %v, %v\n", C.GoString(cformat), C.GoString(arg), err)
	}
	return C.CString(result)
}

// FormatNothing takes a format string and no arguments, and formats it.
//export FormatNothing
func FormatNothing(cformat *C.char) *C.char {
	result, err := pyfmt.Fmt(C.GoString(cformat))
	if err != nil {
		fmt.Printf("Error formatting: %v\n", err)
	}
	return C.CString(result)
}

// FormatOneIntError takes a format string and a single int, and formats it.
//export FormatOneIntError
func FormatOneIntError(cformat *C.char, arg C.int) C.int {
	_, err := pyfmt.Fmt(C.GoString(cformat), int32(arg))
	if err != nil {
		return 1
	}
	return 0
}

// FormatOneFloatError takes a format string and a single float, and formats it.
//export FormatOneFloatError
func FormatOneFloatError(cformat *C.char, arg C.float) C.int {
	_, err := pyfmt.Fmt(C.GoString(cformat), float32(arg))
	if err != nil {
		return 1
	}
	return 0
}

// FormatOneDoubleError takes a format string and a single double, and formats it.
//export FormatOneDoubleError
func FormatOneDoubleError(cformat *C.char, arg C.double) C.int {
	_, err := pyfmt.Fmt(C.GoString(cformat), float64(arg))
	if err != nil {
		return 1
	}
	return 0
}

// FormatOneStringError takes a format string and a single string, and formats it.
//export FormatOneStringError
func FormatOneStringError(cformat *C.char, arg *C.char) C.int {
	_, err := pyfmt.Fmt(C.GoString(cformat), C.GoString(arg))
	if err != nil {
		return 1
	}
	return 0
}

// FormatNothingError takes a format string and no arguments, and formats it. Returns true if the
// formatting errored.
//export FormatNothingError
func FormatNothingError(cformat *C.char) C.int {
	_, err := pyfmt.Fmt(C.GoString(cformat))
	if err != nil {
		return 1
	}
	return 0
}

func main() {}

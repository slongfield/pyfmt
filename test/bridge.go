package main

/*
// TODO(slongfield): Define some types here.
*/
import "C"

import (
	"fmt"

	"github.com/slongfield/pyfmt"
)

// FormatOneInt takes a format string and a single int, and formats it.
//export FormatOneInt
func FormatOneInt(cformat *C.char, arg C.int) *C.char {
	format := C.GoString(cformat)
	result, err := pyfmt.Fmt(format, int32(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v", err)
	}
	return C.CString(result)
}

//export FormatOneFloat
func FormatOneFloat(cformat *C.char, arg C.float) *C.char {
	format := C.GoString(cformat)
	result, err := pyfmt.Fmt(format, float32(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v", err)
	}
	return C.CString(result)
}

//export FormatOneDouble
func FormatOneDouble(cformat *C.char, arg C.double) *C.char {
	format := C.GoString(cformat)
	result, err := pyfmt.Fmt(format, float64(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v", err)
	}
	return C.CString(result)
}

func main() {}

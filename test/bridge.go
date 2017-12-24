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
	result, err := pyfmt.Format(format, int32(arg))
	if err != nil {
		fmt.Printf("Error formatting: %v", err)
	}
	// Note: This allocates memory, and isn't known to the Golang memory manager, so will likely end
	// up leaking.
	// TODO(slongfield): Don't leak memory across these interfaces.
	return C.CString(result)
}

func main() {}

package pyfmt

import (
	"errors"
	"strings"
	"sync"
)

// Using a simple []byte instead of bytes.Buffer to avoid the dependency.
type buffer []byte

func (b *buffer) WriteString(s string) {
	*b = append(*b, s...)
}

const (
	noAlign = iota
	left
	right
	padSign
	center
)

func (b *buffer) WriteAlignedString(s string, align int, width int64, fillChar rune) {
	length := int64(len(s))
	if length >= width || align == noAlign {
		b.WriteString(s)
		return
	}
	var fill string
	if fillChar == 0 {
		fill = " "
	} else {
		fill = string(fillChar)
	}
	switch align {
	case right:
		b.WriteString(strings.Repeat(fill, int(width-length)))
		b.WriteString(s)
	case left:
		b.WriteString(s)
		b.WriteString(strings.Repeat(fill, int(width-length)))
	case center:
		prePad := (width - length) / 2
		b.WriteString(strings.Repeat(fill, int(prePad)))
		b.WriteString(s)
		b.WriteString(strings.Repeat(fill, int(width-length-prePad)))
	// TODO(slongfield): padSign is only valid if we had formatted a number.
	case padSign:
		if s[0] == '-' || s[0] == '+' {
			b.WriteString(string(s[0]))
			b.WriteAlignedString(s[1:], right, width-1, fillChar)
		} else {
			b.WriteAlignedString(s, right, width, fillChar)
		}
	}
}

const (
	useMap = iota
	useList
	useStruct
)

// ff is used to store a formatter's state and is reused iwth sync.Pool to avoid allocations.
type ff struct {
	buf buffer

	// args is the list of arguments passed to Fmt.
	args    []interface{}
	listPos int

	// render renders format parameters
	r render
}

var ffFree = sync.Pool{
	New: func() interface{} { return new(ff) },
}

// newFormater creates a new ff struct.
func newFormater() *ff {
	f := ffFree.Get().(*ff)
	f.listPos = 0
	f.r.init(&f.buf)
	return f
}

func (f *ff) free() {
	f.buf = f.buf[:0]
	f.args = f.args[:0]
	f.listPos = 0
	ffFree.Put(f)
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
			continue
		}
		cachei = i
		for i < end && format[i] != '}' {
			i++
		}
		if format[i] != '}' {
			return errors.New("Single '{' encountered in format string")
		}
		field := format[cachei:i]
		var err error
		name, format := splitFormat(field)
		f.r.val, err = f.getArg(name)
		if err != nil {
			return err
		}
		if formatter, ok := f.r.val.(PyFormatter); ok {
			formatted, err := formatter.PyFormat(format)
			if err != nil {
				return err
			}
			f.buf.WriteString(formatted)
		} else {
			f.r.clearFlags()
			if err = f.r.parseFlags(format); err != nil {
				return err
			}
			if err = f.r.render(); err != nil {
				return err
			}
		}
		i++
	}
	return nil
}

func splitFormat(field string) (string, string) {
	s := strings.SplitN(field, ":", 2)
	if len(s) == 1 {
		return s[0], ""
	}
	return s[0], s[1]
}

func (f *ff) getArg(argName string) (interface{}, error) {
	val, err := getElement(argName, f.listPos, f.args...)
	if argName == "" {
		f.listPos++
	}
	return val, err
}

// Fmt is the equivlent of Python's string.format() function. Takes a list of possible elements
// to use in formatting, and substitutes them.
func Fmt(format string, a ...interface{}) (string, error) {
	f := newFormater()
	defer f.free()
	f.args = a
	err := f.doFormat(format)
	if err != nil {
		return "", err
	}
	s := string(f.buf)
	return s, nil
}

// Must is like Fmt, but panics on error.
func Must(format string, a ...interface{}) string {
	s, err := Fmt(format, a...)
	if err != nil {
		panic(err)
	}
	return s
}

// Error is like Fmt, but returns an error.
func Error(format string, a ...interface{}) error {
	s, err := Fmt(format, a...)
	if err != nil {
		return Error("error formatting {}: {}", s, err)
	}
	return errors.New(s)
}

// PyFormatter is an interface implemented with a PyFormat method that allows for a custom
// formatter.
// The PyFormat method is used to process a custom format spec and then create a formatted version
// of the type based on that.
type PyFormatter interface {
	PyFormat(f string) (string, error)
}

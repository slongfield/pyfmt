package pyfmt

import (
	"errors"
	"sync"
	"unicode/utf8"
)

// buffer type uses a simple []byte instead of bytes.Buffer to avoid the dependency, and has a
// staging buffer which can be used as temporary space to avoid allocations.
type buffer struct {
	contents []byte
	stage    []byte
}

// Implements the io.Writer interface
func (b *buffer) Write(p []byte) (n int, err error) {
	b.contents = append(b.contents, p...)
	return len(p), nil
}

// WriteString writes a string into the backing buffer.
func (b *buffer) WriteString(s string) {
	b.contents = append(b.contents, s...)
}

// WriteString writes a string into the backing buffer 'rep' times.
func (b *buffer) WriteRepeatedString(r string, rep int) {
	b.stage = append(b.stage, r...)
	for len(b.stage) < len(r)*rep {
		b.stage = append(b.stage, b.stage...)
	}
	b.contents = append(b.contents, b.stage[:len(r)*rep]...)
	b.stage = b.stage[:0]
}

const (
	right = iota
	left
	padSign
	center
)

// WriteString writes a string into the backing buffer, padded out to width, based on the alignment
// type.
func (b *buffer) WriteAlignedString(s string, align int, width int64, fillChar rune) {
	length := int64(len(s))
	if length >= width {
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
		b.WriteRepeatedString(fill, int(width-length))
		b.WriteString(s)
	case left:
		b.WriteString(s)
		b.WriteRepeatedString(fill, int(width-length))
	case center:
		prePad := (width - length) / 2
		b.WriteRepeatedString(fill, int(prePad))
		b.WriteString(s)
		b.WriteRepeatedString(fill, int(width-length-prePad))
	case padSign:
		if s[0] == '-' || s[0] == '+' {
			b.WriteString(string(s[0]))
			b.WriteAlignedString(s[1:], right, width-1, fillChar)
		} else {
			b.WriteAlignedString(s, right, width, fillChar)
		}
	}
}

// What type of numbering is being used to access fields. {} is automatic, {0} is manual.
type numbering int
const (
	unknown numbering = iota
	automatic
	manual
)

// ff is used to store a formatter's state and is reused with sync.Pool to avoid allocations.
type ff struct {
	buf buffer

	// args is the list of arguments passed to Fmt.
	args    []interface{}
	listPos int
	numb numbering

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
	f.numb = unknown
	f.r.init(&f.buf)
	return f
}

func (f *ff) free() {
	f.buf.contents = f.buf.contents[:0]
	f.args = f.args[:0]
	f.listPos = 0
	f.numb = unknown
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
				}
				f.buf.WriteString(format[cachei:i])
				i++
				cachei = i
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
		if i < end && format[i] == '{' {
			f.buf.WriteString("{")
			i++
			continue
		}
		cachei = i
		for i < end && format[i] != '}' {
			i++
		}
		if i >= end || format[i] != '}' {
			return errors.New("Single '{' encountered in format string")
		}
		field := format[cachei:i]
		var err error
		name, format := split(field, ':')
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

// Split splits a string on a rune, returning slices pointing to the half before that rune, and
// after. If the rune doesn't appear, the first string returned is the whole string, and the second
// string is empty.
func split(s string, sep rune) (string, string) {
	for i, c := range s {
		if c == sep {
			if i+utf8.RuneLen(sep) <= len(s) {
				return s[:i], s[i+utf8.RuneLen(sep):]
			}
		}
	}
	return s[:], s[len(s):]
}

func (f *ff) getArg(argName string) (interface{}, error) {
	if f.numb == unknown {
		if argName == "" {
			f.numb = automatic
		} else {
			f.numb = manual
		}
	} else {
		if argName == "" && f.numb == manual {
			return nil, Error("cannot switch from manual field specification to automatic field numbering")
		}
		if argName != "" && f.numb == automatic {
			return nil, Error("cannot switch from automatic field numbering to manual field specification")
		}
	} 
	val, err := getElement(argName, f.listPos, f.args...)
	if argName == "" {
		f.listPos++
	}
	return val, err
}

// Fmt is the equivalent of Python's string.format() function. Takes a list of possible elements
// to use in formatting, and substitutes them.
func Fmt(format string, a ...interface{}) (string, error) {
	f := newFormater()
	defer f.free()
	f.args = a
	err := f.doFormat(format)
	if err != nil {
		return "", err
	}
	s := string(f.buf.contents)
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

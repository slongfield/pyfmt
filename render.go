package pyfmt

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"
)

const (
	left = iota
	right
	padSign
	center
)

const (
	signPos = iota
	signNeg
	signSpace
)

const (
	decimal = iota
	binary
	octal
	hex
	hexCap
)

const (
	gen = iota
	genCap
	sci
	sciCap
	fix
	fixCap
	percent
)

type flags struct {
	fillChar   rune
	align      int
	sign       int
	showRadix  bool
	minWidth   uint64
	precision  uint64
	renderType int
}

// Render is the renderer used to render dispatched format strings into a buffer that's been set up
// beforehand.
type render struct {
	buf *buffer
	val interface{}

	flags
}

func (r *render) init(buf *buffer) {
	r.buf = buf
	r.clearFlags()
}

func (r *render) clearFlags() {
	r.flags = flags{}
}

var flagPattern = regexp.MustCompile(`\A(.[<>=^]|[<>=^]?)([\+\- ]?)(#?)(\d*)\.?(\d*)([bdoxXeEfFgG%]?)\z`)

func (r *render) parseFlags(flags string) error {
	if flags == "" {
		return nil
	}
	if !flagPattern.MatchString(flags) {
		// TODO(slongfield): Replace with pyfmt.Error.
		return fmt.Errorf("Invalid flag pattern: %v", flags)
	}
	f := flagPattern.FindStringSubmatch(flags)
	if len(f[1]) > 1 {
		var size int
		r.fillChar, size = utf8.DecodeRuneInString(f[1])
		f[1] = f[1][size:]
	}
	if f[1] != "" {
		switch f[1] {
		case "<":
			r.align = left
		case ">":
			r.align = right
		case "=":
			r.align = padSign
		case "^":
			r.align = center
		default:
			panic("Unreachable, this should never happen.")
		}
	}
	if f[2] != "" {
		switch f[2] {
		case "+":
			r.sign = signPos
		case "-":
			r.sign = signNeg
		case " ":
			r.sign = signSpace
		}
	}
	if f[3] == "#" {
		r.showRadix = true
	}
	if f[4] != "" {
		r.minWidth, _ = strconv.ParseUint(f[4], 10, 64)
	}
	if f[5] != "" {
		r.precision, _ = strconv.ParseUint(f[5], 10, 64)
	}
	if f[6] != "" {
		switch f[6] {
		case "b":
			r.renderType = binary
		case "d":
			r.renderType = decimal
		case "o":
			r.renderType = octal
		case "x":
			r.renderType = hex
		case "X":
			r.renderType = hexCap
		case "e":
			r.renderType = sci
		case "E":
			r.renderType = sciCap
		case "f":
			r.renderType = fix
		case "F":
			r.renderType = fixCap
		case "g":
			r.renderType = gen
		case "G":
			r.renderType = genCap
		case "%":
			r.renderType = percent
		default:
			panic(Must("Unrechable. Saw type match {} not in regex.", f[6]))
		}
	}
	return nil
}

func (r *render) render() error {
	// TODO(slongfield) Dispatch to different functions if the type is set.
	switch t := r.val.(type) {
	case string:
		r.buf.WriteString(r.val.(string))
		return nil
	case int:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int8:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int16:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int32:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int64:
		r.buf.WriteString(strconv.FormatInt(t, 10))
		return nil
	case uint:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint8:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint16:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint32:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint64:
		r.buf.WriteString(strconv.FormatUint(t, 10))
		return nil
	case reflect.Value:
		if t.IsValid() && t.CanInterface() {
			r.val = t.Interface()
			return r.render()
		}
		return r.renderValue(t)
	default:
		r.buf.WriteString(fmt.Sprintf("%v", r.val))
		return nil
	}
}

func (r *render) renderValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		return fmt.Errorf("Invalid value: %v", v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r.buf.WriteString(strconv.FormatInt(int64(v.Int()), 10))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r.buf.WriteString(strconv.FormatUint(uint64(v.Uint()), 10))
		return nil
	case reflect.String:
		r.buf.WriteString(v.String())
		return nil
	default:
		return fmt.Errorf("Unimplemented reflect type %v for %v ", v.Kind(), v)
	}
}

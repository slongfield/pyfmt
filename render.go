package pyfmt

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type flags struct {
	fillChar   rune
	align      int
	sign       string
	showRadix  bool
	minWidth   string
	precision  string
	renderVerb string
	percent    bool
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

var flagPattern = regexp.MustCompile(`\A(.[<>=^]|[<>=^]?)([\+\- ]?)(#?)(\d*)\.?(\d*)([bdoxXeEfFgGrts%]?)\z`)

func (r *render) parseFlags(flags string) error {
	r.renderVerb = "v"
	if flags == "" {
		return nil
	}
	if !flagPattern.MatchString(flags) {
		return Error("Invalid flag pattern: {}", flags)
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
		// "-" is the default behavior, ignore it.
		if f[2] != "-" {
			r.sign = f[2]
		}
	}
	if f[3] == "#" {
		r.showRadix = true
	}
	if f[4] != "" {
		r.minWidth = f[4]
	}
	if f[5] != "" {
		r.precision = "." + f[5]
	}
	if f[6] != "" {
		switch f[6] {
		case "b", "d", "o", "x", "X", "e", "E", "f", "F", "g", "G":
			r.renderVerb = f[6]
		case "%":
			r.percent = true
			r.renderVerb = "f"
		case "r":
			r.renderVerb = "#v"
		case "t":
			r.renderVerb = "T"
		case "s":
			r.renderVerb = "+v"
		default:
			panic("Unreachable, this should never happen. Flag parsing regex is corrupted.")
		}
	}
	return nil
}

func (r *render) render() error {
	var prefix, radix string
	var width int64
	var err error
	if r.percent {
		// Increase the precision by two, to make sure we have enough digits.
		if r.precision == "" {
			r.precision = ".8"
		} else {
			precision, err := strconv.ParseInt(r.precision[1:], 10, 64)
			if err != nil {
				return err
			}
			r.precision = Must(".{}", precision+2)
		}
	}
	if r.showRadix {
		if r.renderVerb == "x" || r.renderVerb == "X" {
			radix = "#"
		} else if r.renderVerb == "b" {
			prefix = "0b"
		} else if r.renderVerb == "o" {
			prefix = "0o"
		}
	}
	if r.align != noAlign {
		width, err = strconv.ParseInt(r.minWidth, 10, 64)
		if err != nil {
			return Error("Can't convert width {} to int", r.minWidth)
		}
		r.minWidth = ""
	}
	str := fmt.Sprintf(strings.Join([]string{
		"%", r.sign, radix, r.minWidth, r.precision, r.renderVerb}, ""), r.val)
	// TODO(slongfield): Add an assertion here that we're operating on a numeric type.
	if prefix != "" {
		if str[0] == '-' {
			str = strings.Join([]string{"-", prefix, str[1:]}, "")
		} else if str[0] == '+' {
			str = strings.Join([]string{"+", prefix, str[1:]}, "")
		} else {
			str = strings.Join([]string{prefix, str}, "")
		}
	}
	// TODO(slongfield): Refactor--pull the percent formatting out and test it
	// independently.
	if r.percent {
		parts := strings.SplitN(str, ".", 2)
		var sign string
		var suffix string
		if len(parts) == 2 {
			prefix, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return Error("Couldn't parse format prefix from: {}", str)
			}
			if prefix == 0 {
				if parts[1][2:] != "" {
					suffix = "." + parts[1][2:]
				}
				if parts[1][0] == '0' {
					str = strings.Join([]string{sign, parts[1][1:2], suffix, "%"}, "")
				} else {
					str = strings.Join([]string{sign, parts[1][0:2], suffix, "%"}, "")
				}
			} else if len(parts[0]) == 1 {
				if parts[1][2:] != "" {
					suffix = "." + parts[1][2:]
				}
				str = strings.Join([]string{sign, parts[0], parts[1][0:2], suffix, "%"}, "")
			} else {
				if parts[1][2:] != "" {
					suffix = "." + parts[1][2:]
				}
				str = strings.Join([]string{sign, parts[0], parts[1][0:2], suffix, "%"}, "")
			}
		} else {
			if _, err := strconv.ParseInt(str, 10, 64); err != nil {
				str = str + "%"
			} else {
				str = str + "00%"
			}
		}
	}
	r.buf.WriteAlignedString(str, r.align, width, r.fillChar)
	return nil
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

package pyfmt

import (
	"errors"
	"fmt"
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
	empty      bool
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

// Flag state machine
const (
	alignState = iota
	signState
	radixState
	zeroState
	widthState
	precisionState
	verbState
	endState
)

// validFlags are 'bdoxXeEfFgGrts%'
func validFlag(b byte) bool {
	return (b == 'b' || b == 'd' || b == 'o' || b == 'x' || b == 'X' || b == 'e' || b == 'E' || b == 'f' || b == 'F' || b == 'g' || b == 'G' || b == 'r' || b == 't' || b == 's' || b == '%')
}

func isDigit(d byte) bool {
	return (d >= '0' && d <= '9')
}

// splitFlags splits out the flags into the various fields.
func splitFlags(flags string) (align, sign, radix, zeroPad, minWidth, precision, verb string, err error) {
	end := len(flags)
	if end == 0 {
		return
	}
	state := alignState
	for i := 0; i < end; {
		switch state {
		case alignState:
			if flags[i] == '<' || flags[i] == '>' || flags[i] == '=' || flags[i] == '^' {
				i = 1
			}
			if end > 1 {
				_, size := utf8.DecodeRuneInString(flags)
				if flags[size] == '<' || flags[size] == '>' || flags[size] == '=' || flags[size] == '^' {
					i = size + 1
				}
			}
			align = flags[0:i]
			state = signState
		case signState:
			if flags[i] == '+' || flags[i] == '-' || flags[i] == ' ' {
				sign = flags[i : i+1]
				i++
			}
			state = radixState
		case radixState:
			if flags[i] == '#' {
				radix = flags[i : i+1]
				i++
			}
			state = zeroState
		case zeroState:
			if flags[i] == '0' {
				zeroPad = flags[i : i+1]
				i++
			}
			state = widthState
		case widthState:
			var j int
			for j = i; j < end; {
				if isDigit(flags[j]) {
					j++
				} else {
					break
				}
			}
			minWidth = flags[i:j]
			i = j
			state = precisionState
		case precisionState:
			if flags[i] == '.' {
				var j int
				for j = i + 1; j < end; {
					if isDigit(flags[j]) {
						j++
					} else {
						break
					}
				}
				precision = flags[i:j]
				i = j
			}
			state = verbState
		case verbState:
			if validFlag(flags[i]) {
				verb = flags[i : i+1]
				i++
			}
			state = endState
		default:
			// Get to this state when we've run out of other states. If we reach this, it means we've
			// gone too far, since we've passed the verb state, but aren't at the end of the string, so
			// error.
			err = errors.New("Could not decode format specification: " + flags)
			i = end + 1
		}
	}
	return
}

func (r *render) parseFlags(flags string) error {
	r.renderVerb = "v"
	if flags == "" {
		r.empty = true
		return nil
	}
	align, sign, radix, zeroPad, minWidth, precision, verb, err := splitFlags(flags)
	if err != nil {
		return Error("Invalid flag pattern: {}, {}", flags, err)
	}
	if len(align) > 1 {
		var size int
		r.fillChar, size = utf8.DecodeRuneInString(align)
		align = align[size:]
	}
	if align != "" {
		switch align {
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
	if sign != "" {
		// "-" is the default behavior, ignore it.
		if sign != "-" {
			r.sign = sign
		}
	}
	if radix == "#" {
		r.showRadix = true
	}
	if zeroPad != "" {
		if align == "" {
			r.align = padSign
		}
		if r.fillChar == 0 {
			r.fillChar = '0'
		}
	}
	if minWidth != "" {
		r.minWidth = minWidth
	}
	if precision != "" {
		r.precision = precision
	}
	if verb != "" {
		switch verb {
		case "b", "o", "x", "X", "e", "E", "f", "F", "g", "G":
			r.renderVerb = verb
		case "d":
			r.renderVerb = verb
			r.showRadix = false
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

// render renders a single element by passing that element and the translated format string
// into the fmt formatter.
func (r *render) render() error {
	var prefix, radix string
	var width int64
	var err error
	if r.empty {
		fmt.Fprint(r.buf, r.val)
		return nil
	}

	if r.percent {
		if err = r.setupPercent(); err != nil {
			return err
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

	if r.minWidth == "" {
		width = 0
	} else {
		width, err = strconv.ParseInt(r.minWidth, 10, 64)
		if err != nil {
			return Error("Can't convert width {} to int", r.minWidth)
		}
	}

	// Only let Go handle the width for floating+complex types, elsewhere the alignment rules are
	// different.
	if r.renderVerb != "f" && r.renderVerb != "F" && r.renderVerb != "g" && r.renderVerb != "G" && r.renderVerb != "e" && r.renderVerb != "E" {
		r.minWidth = ""
	}

	str := fmt.Sprintf("%"+r.sign+radix+r.minWidth+r.precision+r.renderVerb, r.val)

	if prefix != "" {
		// Get rid of any prefix added by minWidth. We'll add this back in later when we
		// WriteAlignedString to the underlying buffer
		str = strings.TrimLeft(str, " ")
		if str != string(r.fillChar) {
			str = strings.TrimLeft(str, string(r.fillChar))
		}
		if len(str) > 0 && str[0] == '-' {
			str = strings.Join([]string{"-", prefix, str[1:]}, "")
		} else if len(str) > 0 && str[0] == '+' {
			str = strings.Join([]string{"+", prefix, str[1:]}, "")
		} else if r.sign == " " {
			str = strings.Join([]string{" ", prefix, str}, "")
		} else {
			str = strings.Join([]string{prefix, str}, "")
		}
	}

	if r.renderVerb == "f" || r.renderVerb == "F" || r.renderVerb == "g" || r.renderVerb == "G" || r.renderVerb == "e" || r.renderVerb == "E" {
		str = strings.TrimSpace(str)
		if r.sign == " " && str[0] != '-' {
			str = " " + str
		}
	}

	if r.percent {
		str, err = transformPercent(str)
		if err != nil {
			return err
		}
	}

	if len(str) > 0 {
		if str[0] != '(' && (r.align == left || r.align == padSign) {
			if str[0] == '-' {
				r.buf.WriteString("-")
				str = str[1:]
				width--
			} else if str[0] == '+' {
				r.buf.WriteString("+")
				str = str[1:]
				width--
			} else if str[0] == ' ' {
				r.buf.WriteString(" ")
				str = str[1:]
				width--
			} else {
				r.buf.WriteString(r.sign)
			}
		}
	}

	if r.showRadix && r.align == padSign {
		r.buf.WriteString(str[0:2])
		r.buf.WriteAlignedString(str[2:], r.align, width-2, r.fillChar)
	} else {
		r.buf.WriteAlignedString(str, r.align, width, r.fillChar)
	}
	return nil
}

func (r *render) setupPercent() error {
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
	return nil
}

func transformPercent(p string) (string, error) {
	var sign string
	if p[0] == '-' {
		sign = "-"
		p = p[1:]
	}
	intPart, mantissa := split(p, '.')
	var suffix string
	if mantissa != "" {
		prefix, err := strconv.ParseInt(intPart, 10, 64)
		if err != nil {
			return "", Error("Couldn't parse format prefix from: {}", p)
		}
		if prefix == 0 {
			if mantissa[2:] != "" {
				suffix = "." + mantissa[2:]
			}
			if mantissa[0] == '0' {
				return strings.Join([]string{sign, mantissa[1:2], suffix, "%"}, ""), nil
			}
			return strings.Join([]string{sign, mantissa[0:2], suffix, "%"}, ""), nil
		} else if len(intPart) == 1 {
			if mantissa[2:] != "" {
				suffix = "." + mantissa[2:]
			}
			return strings.Join([]string{sign, intPart, mantissa[0:2], suffix, "%"}, ""), nil
		}
		if mantissa[2:] != "" {
			suffix = "." + mantissa[2:]
		}
		return strings.Join([]string{sign, intPart, mantissa[0:2], suffix, "%"}, ""), nil
	}
	if _, err := strconv.ParseInt(p, 10, 64); err != nil {
		return p + "%", nil
	}
	return p + "00%", nil

}

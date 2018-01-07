package pyfmt

import (
	"fmt"
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

var flagPattern = regexp.MustCompile(`\A((?:.[<>=^])|(?:[<>=^])?)([\+\- ]?)(#?)(0?)(\d*)(\.\d*)?([bdoxXeEfFgGrts%]?)\z`)

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
		if f[1] == "" {
			r.align = padSign
		}
		if r.fillChar == 0 {
			r.fillChar = '0'
		}
	}
	if f[5] != "" {
		r.minWidth = f[5]
	}
	if f[6] != "" {
		r.precision = f[6]
	}
	if f[7] != "" {
		switch f[7] {
		case "b", "o", "x", "X", "e", "E", "f", "F", "g", "G":
			r.renderVerb = f[7]
		case "d":
			r.renderVerb = f[7]
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

	str := fmt.Sprintf(strings.Join([]string{
		"%", r.sign, radix, r.minWidth, r.precision, r.renderVerb}, ""), r.val)

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
				width -= 1
			} else if str[0] == '+' {
				r.buf.WriteString("+")
				str = str[1:]
				width -= 1
			} else if str[0] == ' ' {
				r.buf.WriteString(" ")
				str = str[1:]
				width -= 1
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
	var parts []string
	var sign string
	if p[0] == '-' {
		sign = "-"
		p = p[1:]
	}
	parts = strings.SplitN(p, ".", 2)
	var suffix string
	if len(parts) == 2 {
		prefix, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return "", Error("Couldn't parse format prefix from: {}", p)
		}
		if prefix == 0 {
			if parts[1][2:] != "" {
				suffix = "." + parts[1][2:]
			}
			if parts[1][0] == '0' {
				return strings.Join([]string{sign, parts[1][1:2], suffix, "%"}, ""), nil
			} else {
				return strings.Join([]string{sign, parts[1][0:2], suffix, "%"}, ""), nil
			}
		} else if len(parts[0]) == 1 {
			if parts[1][2:] != "" {
				suffix = "." + parts[1][2:]
			}
			return strings.Join([]string{sign, parts[0], parts[1][0:2], suffix, "%"}, ""), nil
		}
		if parts[1][2:] != "" {
			suffix = "." + parts[1][2:]
		}
		return strings.Join([]string{sign, parts[0], parts[1][0:2], suffix, "%"}, ""), nil
	}
	if _, err := strconv.ParseInt(p, 10, 64); err != nil {
		return p + "%", nil
	}
	return p + "00%", nil

}

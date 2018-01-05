// Copyright header goes here.

/*
pyfmt

pyfmt implements Python-style advanced string formatting.

This is an alternative to the `fmt` package's Sprintf-style string formatting, and mimics the
.format() style formatting available in Python >2.6.

Braces {} are used to indicate 'format items', anything outside of braces will be emitted directly
into the output string, and anything inside will be used to get values from the other function
arguments, and format them. The one exception is double braces, '{{' and '}}', which will cause
literal '{' and '}' runes to be emitted in the output.

Each format item consists of a 'field name', which indicates which value from the argument list to
use, and a 'format specifier', which indicates how to format that item.

Functions

pyfmt implements three functions, 'Fmt', 'Must', and 'Error'. 'Fmt' formats, but may return
an error as detailed below. 'Must' formats, but will panic when 'Fmt' would return an error, and
'Error' acts like 'Fmt', but returns an error type. In the event that there's an error formatting
the error, 'Error' includes the format error and as much of the formatted string as possible.

All of them take a format string, and then a list of arguments to look up elements from.

Getting Values from Field Names

Values can be fetched from field names in two forms: simple names, or compound names. All compound
names build off of simple names, and all simple names are dependent on the type of format function
you call.

Simple field names:

The simplest look up treats the argument list as just a list. There are two possible ways to look
up elements from this list. First, by {}, which gets the 'next' item, and second, by {n}, which
gets the nth item. Accessing these two ways is independent, so while

  pyfmt.Must("{} {} {}", ...)

Is equivalent to

  pyfmt.Must("{0} {1} {2}", ...)

so is

  pyfmt.Must("{} {1} {2}", ...)

but

  pyfmt.Must("{} {1} {}", ...)

Is equivalent to:

  pyfmt.Must("{0} {1} {1}". ...)

Accessing an element that's outside the list range will return an error (with Fmt) or panic (with
Must).

The first element in the list is treated specially if it's a struct or a map with string keys,
allowing the elements from that struct or map can be directly accessed. For instance:

  pyfmt.Must("{test}", map[string]int{"test": 5})

returns

  "5".

and for structs:

  pyfmt.Must("{test}": myStruct{test: 5})

returns

 "5".

Attempting to read from an undefined key will return an error or panic, depending on if it was
accessed with Fmt or Must..

Compound field names:

If the value referenced by the field is itself a List, map[string]interface{}, or struct, it can be
further accessed in the format string.

Lists are accessed with square brackets:

  pyfmt.Must("{0[0]}", []string{"test"}) -> "test"

Similarly, maps are accessed with square brackets:

  pyfmt.Must("{0[test]}", map[string]interface{}{"test": "42"}) -> "42"

And struct fields are accessed with period, '.'

  pyfmt.Must("{foo.bar.baz}", MyStruct{foo: Foo{bar: Bar{baz: "test"}}}) -> "test"

Formatting

If after a field name, there's a ':', what follows is considered to be the format specifier. If a
type satisfies the Format interface (discussed below), the format specifier will be passed to that,
but otherwise, it will fall back to the default formatter, which expects the standard format
specifier:

  [[fill]align][sign][#][0][minimumwidth][.precision][type]

The optional align feature can be one of the following:

  '<': left-aligned
  '>': right-aligned
  '=': padding after the sign, but before the digits (e.g., for +000042)
  '^': centered

If an align flag is defined, a 'fill' character can also be defined. If undefined, space (' ') will
be used.

The optional 'sign' is only valid for numeric types and can be:

  '+': Show sign for both positive and negative numbers
  '-': Show sign only for negative numbers (default)
  ' ': use a leading space for positive numbers

If # is present, when using the binary, octal, or hex types, a '0b', '0o', or '0x' will be
prepended, respectively.

The minimumwidth field specifies a minimum width, which is helpful when used with alignment. If
preceded with a zero, numbers will be zero-padded.

The precision field specifies a maximum width for non-floating point, non-integer types, and the
number of points to show after the decimal point for floating types.

The 'type' format determines what type the value will be formatted as.

For integers:

  'b' - Binary, base 2
  'd' - Decimal, base 10 (default)
  'o' - Octal, base 8
  'x' - Hexadecimal, base 16
  'X' - Hexadecimal, base 16, using upper-case letters

For floats and complex numbers:

  'e' - Scientific notation
  'E' - Similar to e, but uppercase
  'f' - Fixed point, displays number as a fixed-point number.
  'F' - Same, but uppercase.
  'g' - General format, prints as a fixed point unless its too large, than switches to scientific
        notation. (default)
  'G' - Similar to g, but uses capital letters
  '%' - Percentage, multiplies the number by 100 and displays it with a '%' sign. Can also be
        applied to integer types.

Special Formatting Types

For some types (most notably structs), the default formatter doesn't quite give enough
information to understand the value after its printed, so it's useful to get more accurate Go
representations. Additionally, sometimes it's useful to print the type of a variable while
formatting it. For these, pyfmt allows for some special formatting types that aren't in the
Python format syntax.

  'r' - convert the value to its Go-syntax representation
  't' - convert the value to its Go type
  's' - if printing a struct, print the struct field names

These are equivalent to the `%#v`, `%T` and `%+v` format strings in the "fmt" package.

Custom formatters

Internally, pyfmt uses Go's fmt package, so existing types satisfying those Formatter, GoStringer,
or Stringer interfaces will use those implementations as appropriate.

Examples

TODO(slongfield): Add examples.

TODO(slongfield): Details on custom formatters.

TODOs

  *  Add support for alignment
  *  Add support for custom formatters
  *  Add more tests.
*/
package pyfmt

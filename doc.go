// Copyright header goes here.

/*
pyfmt implements Python-style advanced string formatting.

This is an alternative to the `fmt` package's Sprintf-style string formatting, and mimics the
.format() style formatting available in Python >2.6.

Braces {} are used to indicate 'format items', anything outside of braces will be emitted directly
into the output string, and anything inside will be used to get values from the other function
arguments, and format them. The one exception is double braces, '{{' and '}}', which will cause
literal '{' and '}' runes to be emitted in the output.

Each format item consists of a 'field name', which indicates which value from the argument list to
use, and a 'format specifier', which indicates how to format that item.

Getting Values from Field Names

Values can be fetched from field names in two forms: simple names, or compound names. All compound
names build off of simple names, and all simple names are dependent on the type of format function
you call.

Simple field names:

Using the base Format function, you can lookup fields in two ways from lists, first, by {}, which
gets the 'next' item, and second, by {n}, which gets the nth item. Accessing these two ways is
independent, so while

  Format("{} {} {}", ...)

Is equivalent to

  Format("{0} {1} {2}", ...)

so is

  Format("{} {1} {2}", ...)

but

  Format("{} {1} {}", ...)

Is equivalent to:

  Format("{0} {1} {1}". ...)

Accessing an element that's outside the list range, will return an error, and MustFormat will panic.

Using FormatMap, pass a map[string]interface{} to lookup names from the map using the string keys
from the map. For instance:

  FormatMap("{test}", map[string]int{"test": 5})

returns

  "5".

Attempting to read from an undefined key will return an error.

Using FormatStruct, you can reference arguments from a struct:

  FormatStruct("{test}": myStruct{test: 5})

returns

 "5".

Compound field names:

If the value referenced by the field is itself a List, map[string]interface{}, or struct, it can be further accessed in the format string.

Lists are accessed with square brackets:

  Format("{0[0]}", []string{"test"}) -> "test"

Similarly, maps are accessed with square brackets:

  Format("{0[test]}", map[string]interface{}{"test": "42"}) -> "42"

And struct fields are accessed with period, '.'

  FormatStruct("{foo.bar.baz}", MyStruct{foo: Foo{bar: Bar{baz: "test"}}}) -> "test"

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
preceeded with a zero, numbers will be zero-padded.

The precision field specifies a maximum width for non-floating point, non-integer types, and the
number of points to show after the decimal point for floating types.

The 'type' format determines what type the value will be formatted as.

For integers:

  'b' - Binary, base 2
	'd' - Decimal, base 10 (default)
	'o' - Octal, base 8
	'x' - Hexadecimal, base 16
	'X' - Hexadecimal, base 16, using upper-case letters

For floats:

  'e' - Scientific notation
	'E' - Similar to e, but uppercase
	'f' - Fixed point, displays number as a fixed-point number.
	'F' - Same, but uppercase.
	'g' - General format, prints as a fixed point unless its too large, than switches to scientific
	      notation. (default)
	'G' - Similar to g, but uses capital letters
	'%' - Percentage, multiplies the number by 100 and displays it with a '%' sign. Can also be
	      applied to integer types.

Examples

TODO(slongfield): Add examples.

TODO(slongfield): Details on custom formatters.
*/
package pyfmt

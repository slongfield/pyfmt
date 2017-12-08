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

TODO(slongfield): Make notes about how to get values out of lists, maps, and structs, as well as
notes about nested lists, maps, and structs.

Compound field names:

Formatting

Note note note

Examples

TODO(slongfield): Add examples.
*/
package pyfmt

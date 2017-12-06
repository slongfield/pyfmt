// Copyright header goes here.

/*
pyfmt implements Python-style advanced string formatting.

This is an alternative to the `fmt` package's Sprintf-style string formatting, and duplicates the
.format() style formatting available in Python >2.6.

Braces {} are used to indicate 'fields', anything outside of braces will be emitted directly into
the output string, and anything inside will be used to get values from the other function arguments,
and format them. The one exception is double braces, '{{' and '}}', which will cause literal '{' and
'}' runes to be emitted in the output.

Each field consists of a 'field name', which indicates which value from the argument list to use,

Getting Values from Field Names

TODO(slongfield): Make notes about how to get values out of lists, maps, and structs, as well as
notes about nested lists, maps, and structs.

  Simple field names

  Compound field names

Formatting



Examples

TODO(slongfield): Add examples.
*/
package pyfmt

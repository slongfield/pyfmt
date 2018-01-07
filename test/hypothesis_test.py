"""Use the Python Hypotehsis library to test equivalence.

This uses constrained random generation to generate thousands of different examples of format
strings and then formats them both with Python's format library, and the Go implementation in
pyfmt.

Does not test for error equality, instead tests for:
    If the Python formatter formats it without error, the Go formatter will format it in the same
    way.
"""

import string
from hypothesis import given, assume, settings
from hypothesis.strategies import text, from_regex, integers

import build


@given(text(alphabet=string.printable))
@settings(max_examples=1000)
def test_no_format_arguments(fmt_str):
    """Test that without format arguments, it's OK."""
    # First format it with pyfmt. If it doesn't format correctly, toss out this test run.
    try:
        pyfmt = fmt_str.format().encode("ascii")
    except (ValueError, IndexError, KeyError):
        assume(False)

    gofmt = build.FormatNothing(fmt_str.encode("ascii"))

    assert gofmt == pyfmt


@given(text(alphabet=string.printable, max_size=10),
       from_regex(r"\A(([:print:][<>=^])|([<>=^]?))[\+\- ]?#?[0-9]{0,4}[bdoxX]\Z"),
       text(alphabet=string.printable, max_size=10),
       integers(min_value=-2**31, max_value=2**31 - 1))
@settings(max_examples=20000)
def test_format_one_int(pre_str, fmt, post_str, val):
    """Test that a single integer is formatted correctly."""
    fmt_str = pre_str + "{:" + fmt + "}" + post_str

    try:
        pyfmt = fmt_str.format(val).encode("ascii")
    except (ValueError, IndexError, KeyError):
        assume(False)

    print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneInt(fmt_str.encode("ascii"), val)

    assert gofmt == pyfmt

"""Use the Python Hypotehsis library to test equivalence.

This uses constrained random generation to generate thousands of different examples of format
strings and then formats them both with Python's format library, and the Go implementation in
pyfmt.

Does not test for error equality, instead tests for:
    If the Python formatter formats it without error, the Go formatter will format it in the same
    way.
"""

import string
from hypothesis import assume, given, settings
from hypothesis.strategies import floats, from_regex, integers, text

import build

_NUM_TEST = 10000
_DEBUG = False


@given(text(alphabet=string.printable))
@settings(max_examples=_NUM_TEST)
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
@settings(max_examples=_NUM_TEST)
def test_format_one_int(pre_str, fmt, post_str, val):
    """Test that a single integer is formatted correctly."""
    fmt_str = pre_str + "{:" + fmt + "}" + post_str

    try:
        pyfmt = fmt_str.format(val).encode("ascii")
    except (ValueError, IndexError, KeyError):
        assume(False)

    if _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneInt(fmt_str.encode("ascii"), val)

    assert gofmt == pyfmt


@given(text(alphabet=string.printable, max_size=10),
       from_regex(
           r"\A(([:print:][<>=^])|([<>=^]?))[\+\- ]?0?[1-9][0-9]{0,3}"
           r"\.[1-9][0-9]{0,6}[eEfFgG]\Z"),
       text(alphabet=string.printable, max_size=10),
       floats(allow_nan=False, allow_infinity=False))
@settings(max_examples=_NUM_TEST)
def test_format_one_double(pre_str, fmt, post_str, val):
    """Test that a single double is formatted correctly.

    Note that the format string always specifies the width and precision, since Python and Golang
    have different defaults.

    Does a bit of post-proessing:
        Replaces .E with E -- Python renders 0 as 0.E+00, Go uses OE+00

    Disallowing NaN and infinity. Tested those within go, and there are uninteresting rendering
    differences beteen Golang and Python here.
    """
    fmt_str = pre_str.replace(".E", "E") + "{:" + fmt + "}" + post_str.replace(".E", "E")

    try:
        pyfmt = fmt_str.format(val).replace(".E", "E").encode("ascii")
    except (ValueError, IndexError, KeyError):
        assume(False)

    if len(pyfmt) < 100 and _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneDouble(fmt_str.encode("ascii"), val)

    assert gofmt == pyfmt


@given(text(alphabet=string.printable, max_size=10),
       from_regex(r"\A([:print:][<>=^][0-9]{0,4})?\Z"),
       text(alphabet=string.printable, max_size=10),
       text(alphabet=string.printable, max_size=10))
@settings(max_examples=_NUM_TEST)
def test_format_one_str(pre_str, fmt, post_str, val):
    """Test that a single string is formatted correctly.

    Note that, for this test, whenever alignment is requested, explicitly request the type of
    alignment. pyfmt has a default alignment of right aligned for both strings and numbers, but
    Python has a defalt alignment of right for numbers, and a default alignment of left for
    strings, despite the PEP3101 documentation saying that the alignment is left by default.
    """
    fmt_str = pre_str + "{:" + fmt + "}" + post_str

    try:
        pyfmt = fmt_str.format(val).encode("ascii")
    except (ValueError, IndexError, KeyError):
        assume(False)

    if _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneString(fmt_str.encode("ascii"), val.encode("ascii"))

    assert gofmt == pyfmt

"""Use the Python Hypotehsis library to test equivalence.

This uses constrained random generation to generate thousands of different examples of format
strings and then formats them both with Python's format library, and the Go implementation in
pyfmt.
"""

import string
from hypothesis import assume, given, settings, HealthCheck
from hypothesis.strategies import floats, from_regex, integers, text

import build

_NUM_TEST = 1000
_DEBUG = True


@given(text(alphabet=string.printable))
@settings(max_examples=_NUM_TEST)
def test_no_format_arguments(fmt_str):
    """Test that without format arguments, it's OK."""
    # First format it with pyfmt. If it doesn't format correctly, toss out this test run.
    try:
        pyfmt = fmt_str.format().encode("ascii")
    except (ValueError, IndexError, KeyError):
        assert build.FormatNothingError(fmt_str.encode("ascii"))
        return

    gofmt = build.FormatNothing(fmt_str.encode("ascii"))

    assert gofmt == pyfmt


@given(text(alphabet=string.printable))
@settings(max_examples=_NUM_TEST)
def test_no_format_arguments_errors(fmt_str):
    """Test that without format arguments, match error conditions."""
    if _DEBUG:
        print("{}.format()".format(fmt_str))

    # If python errors, golang should error.
    try:
        fmt_str.format().encode("ascii")
    except (ValueError, IndexError, KeyError):
        assert build.FormatNothingError(fmt_str.encode("ascii"))
        return

    # If python did not error, Golang shouldn't error.
    try:
        build.FormatNothing(fmt_str.encode("ascii"))
    except e:
        assert False


@given(text(alphabet=string.printable, max_size=10),
       from_regex(r"\A(([:print:][<>=^])|([<>=^]?))[\+\- ]?#?[0-9]{0,4}[bdoxX]\Z"),
       text(alphabet=string.printable, max_size=10),
       integers(min_value=-2**31, max_value=2**31 - 1))
@settings(max_examples=_NUM_TEST,suppress_health_check=(HealthCheck.filter_too_much,))
def test_format_one_int(pre_str, fmt, post_str, val):
    """Test that a single integer is formatted correctly."""
    fmt_str = pre_str + "{:" + fmt + "}" + post_str

    try:
        pyfmt = fmt_str.format(val).encode("ascii")
    except (ValueError, IndexError, KeyError) as exp:
        # Ignore the cases where formatting goes from {} to {0} formatting. The Go
        # implementation supports this, but Python's doesn't
        if "cannot switch from automatic field numbering" in str(exp):
            assume(False)

        assert build.FormatOneIntError(fmt_str.encode("ascii"), val)
        return

    if _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneInt(fmt_str.encode("ascii"), val)

    assert gofmt == pyfmt


@given(text(alphabet=string.printable, max_size=10),
       from_regex(r"\A(([:print:][<>=^])|([<>=^]?))[\+\- ]?#?[0-9]{0,4}[bdoxX]\Z"),
       text(alphabet=string.printable, max_size=10),
       integers(min_value=-2**31, max_value=2**31 - 1))
@settings(max_examples=_NUM_TEST,suppress_health_check=(HealthCheck.filter_too_much,))
def test_format_one_int_erros(pre_str, fmt, post_str, val):
    """Test that a single integer is formatted correctly."""
    fmt_str = pre_str + "{:" + fmt + "}" + post_str

    try:
        fmt_str.format(val).encode("ascii")
    except (ValueError, IndexError, KeyError) as exp:
        # Ignore the cases where formatting goes from {} to {0} formatting. The Go
        # implementation supports this, but Python's doesn't
        if "cannot switch from automatic field numbering" in str(exp):
            assume(False)

        assert build.FormatOneIntError(fmt_str.encode("ascii"), val)
        return
    assume(False)


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
    except (ValueError, IndexError, KeyError) as exp:
        if "cannot switch from automatic field numbering" in str(exp):
            assume(False)
        assert build.FormatOneDoubleError(fmt_str.encode("ascii"), val)
        return

    if len(pyfmt) < 100 and _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneDouble(fmt_str.encode("ascii"), val)

    assert gofmt == pyfmt


@given(text(alphabet=string.printable, max_size=10),
       from_regex(r"\A([:print:][<>^][0-9]{0,4})?\Z"),
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
    except (ValueError, IndexError, KeyError) as exp:
        if "cannot switch from automatic field numbering" in str(exp):
            assume(False)
        assert build.FormatOneStringError(fmt_str.encode("ascii"), val.encode("ascii"))
        return

    if _DEBUG:
        print("{}.format({}) = {}".format(fmt_str, val, pyfmt))

    gofmt = build.FormatOneString(fmt_str.encode("ascii"), val.encode("ascii"))

    assert gofmt == pyfmt

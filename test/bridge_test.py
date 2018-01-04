"""Testing the Go implementation of pyfmt through the bridge."""
import math

import build
import pytest


@pytest.mark.parametrize("val", [42, -10, 100000, 0, 2**31 - 1])
@pytest.mark.parametrize("fmt_str", ["{}", "{:b}", "{:x}", "{:d}", "{:X}", "{:#x}", "{:#X}",
                                     "{:#d}", "{:#b}", "{:#o}"])
def test_int(val, fmt_str):
    """Simple tests of integer fomatting."""
    gofmt = build.FormatOneInt(fmt_str.encode("ascii"), val)
    pyfmt = fmt_str.format(val).encode("ascii")
    assert gofmt == pyfmt


# Since default precision is different between Golang and Python, always specify the precision for
# float and double format strings. Also, since Python seems to always use float64 without Numpy,
# some of the test cases don't work quite the same for float32 and float64.
@pytest.mark.parametrize("val", [1.2, -1.2, 3.0 / 4.0, -1, 0, float('nan'), 2.0**20])
@pytest.mark.parametrize("fmt_str", ["{:.6e}", "{:.6E}", "{:.6f}", "{:.6F}",
                                     "{:.6g}", "{:.6G}"])
def test_float(val, fmt_str):
    """Simple tests of float formatting."""
    gofmt = build.FormatOneFloat(fmt_str.encode("ascii"), val)
    pyfmt = fmt_str.format(val).encode("ascii")
    # Go renders nan as NaN, but Python renders as nan. NaN is better.
    if math.isnan(val):
        assert gofmt.lower() == pyfmt.lower()
    else:
        assert gofmt == pyfmt


@pytest.mark.parametrize("val", [1.2, -1.2, 3.0 / 4.0, 1.0 / 11.0, -1, 0, float('nan'), 2.1**20])
@pytest.mark.parametrize("fmt_str", ["{:.6e}", "{:.6E}", "{:.6f}", "{:.6F}", "{:.6g}",
                                     "{:.6G}", "{:5.5f}", "{:+4.4e}", "{:-3.3g}",
                                     "{: 1.7F}"])
def test_double(val, fmt_str):
    """Simple tests of double formatting."""
    gofmt = build.FormatOneDouble(fmt_str.encode("ascii"), val)
    pyfmt = fmt_str.format(val).encode("ascii")
    if math.isnan(val):
        assert gofmt.lower() == pyfmt.lower()
    else:
        assert gofmt == pyfmt

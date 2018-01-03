"""Testing the Go implementation of pyfmt through the bridge."""
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


@pytest.mark.parametrize("val", [1.2, -1.2, 3.0 / 4.0, -1, 0, float('nan')])
@pytest.mark.parametrize("fmt_str", ["{}"])
def test_float(val, fmt_str):
    """Simple tests of float formatting."""
    gofmt = build.FormatOneFloat(fmt_str.encode("ascii"), val)
    pyfmt = fmt_str.format(val).encode("ascii")
    # Go renders nan as NaN, but Python renders as nan. NaN is better.
    assert gofmt.lower() == pyfmt


@pytest.mark.parametrize("val", [1.2, -1.2, 3.0 / 4.0, 1.0 / 11.0, -1, 0, float('nan')])
@pytest.mark.parametrize("fmt_str", ["{}"])
def test_double(val, fmt_str):
    """Simple tests of double formatting."""
    gofmt = build.FormatOneDouble(fmt_str.encode("ascii"), val)
    pyfmt = fmt_str.format(val).encode("ascii")
    assert gofmt.lower() == pyfmt

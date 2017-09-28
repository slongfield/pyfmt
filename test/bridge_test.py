import build


def test_simple():
    assert build.FormatOneInt("%d".encode("ascii"), 42) == "42".encode("ascii")

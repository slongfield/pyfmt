import build


def test_simple():
    assert build.FormatOneInt("{}".encode("ascii"), 42) == "42".encode("ascii")

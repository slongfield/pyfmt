"""py.test benchmark for using the FFI bridge vs compiling Go.

Using the bridge is, unsurprisingly, about 1000 times faster."""

import subprocess
import textwrap
import build


def test_simple(benchmark):
    """Benchmark the bridge."""
    result = benchmark(lambda: build.FormatOneInt("%d".encode("ascii"), 42))
    assert result == "42".encode("ascii")


def test_compile(benchmark):
    """Benchmark compilation."""
    with open("/tmp/test.go", 'w') as tempfile:
        tempfile.write(textwrap.dedent("""
            package main

            import (
              "fmt"

              "github.com/slongfield/pyfmt"
             )

             func main() {
               out, _ := pyfmt.Format("%d", 42)
               fmt.Print(out)
             }
            """))
    out = benchmark(lambda: subprocess.check_output(['go', 'run',
                                                     '/tmp/test.go']).decode('utf-8'))
    assert out == "42"

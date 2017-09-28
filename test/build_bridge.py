#!/usr/bin/python3
from cffi import FFI

ffi = FFI()
ffi.set_source("._bridge", None)

# Copied from the Go-generated libbridge.h
ffi.cdef("""
char* FormatOneInt(char* p0, int p1);
""")

if __name__ == "__main__":
    ffi.compile()

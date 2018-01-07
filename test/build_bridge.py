#!/usr/bin/python3
from cffi import FFI

ffi = FFI()
ffi.set_source("._bridge", None)

# Copied from the Go-generated libbridge.h
ffi.cdef("""
char* FormatOneInt(char* p0, int p1);
char* FormatOneFloat(char* p0, float p1);
char* FormatOneDouble(char* p0, double p1);
char* FormatOneString(char* p0, char* p1);
char* FormatNothing(char* p0);
""")

if __name__ == "__main__":
    ffi.compile()

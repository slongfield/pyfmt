import os
from _bridge import ffi

lib = ffi.dlopen(os.path.join(os.path.dirname(__file__), "libbridge.so"))


def FormatOneInt(fmt, a):
    return ffi.string(lib.FormatOneInt(fmt, a))


def FormatOneFloat(fmt, a):
    return ffi.string(lib.FormatOneFloat(fmt, a))


def FormatOneDouble(fmt, a):
    return ffi.string(lib.FormatOneDouble(fmt, a))


def FormatOneString(fmt, a):
    return ffi.string(lib.FormatOneString(fmt, a))


def FormatNothing(fmt):
    return ffi.string(lib.FormatNothing(fmt))


def FormatOneIntError(fmt, a):
    return lib.FormatOneIntError(fmt, a)


def FormatOneFloatError(fmt, a):
    return lib.FormatOneFloatError(fmt, a)


def FormatOneDoubleError(fmt, a):
    return lib.FormatOneDoubleError(fmt, a)


def FormatOneStringError(fmt, a):
    return lib.FormatOneStringError(fmt, a)


def FormatNothingError(fmt):
    return lib.FormatNothingError(fmt)

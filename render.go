package pyfmt

import (
	"errors"
	"strconv"
)

type flags struct {
}

// Render is the renderer used to render dispatched format strings into a buffer that's been set up
// beforehand.
type render struct {
	buf *buffer
	val interface{}

	flags
}

func (r *render) init(buf *buffer) {
	r.buf = buf
	r.clearflags()
}

func (r *render) clearflags() {
	r.flags = flags{}
}

func (r *render) render() error {
	// TODO(slongfield) Dispatch to different functions if the type is set.
	switch r.val.(type) {
	case string:
		r.buf.WriteString(r.val.(string))
		return nil
	case int:
		r.buf.WriteString(strconv.FormatInt(int64(r.val.(int)), 10))
		return nil
	case int8:
		r.buf.WriteString(strconv.FormatInt(int64(r.val.(int8)), 10))
		return nil
	case int16:
		r.buf.WriteString(strconv.FormatInt(int64(r.val.(int16)), 10))
		return nil
	case int32:
		r.buf.WriteString(strconv.FormatInt(int64(r.val.(int32)), 10))
		return nil
	case int64:
		r.buf.WriteString(strconv.FormatInt(r.val.(int64), 10))
		return nil
	case uint:
		r.buf.WriteString(strconv.FormatUint(uint64(r.val.(uint)), 10))
		return nil
	case uint8:
		r.buf.WriteString(strconv.FormatUint(uint64(r.val.(uint8)), 10))
		return nil
	case uint16:
		r.buf.WriteString(strconv.FormatUint(uint64(r.val.(uint16)), 10))
		return nil
	case uint32:
		r.buf.WriteString(strconv.FormatUint(uint64(r.val.(uint32)), 10))
		return nil
	case uint64:
		r.buf.WriteString(strconv.FormatUint(r.val.(uint64), 10))
		return nil

	default:
		return errors.New("Unimplemented!")
	}
}

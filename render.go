package pyfmt

import (
	"fmt"
	"reflect"
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
	r.clearFlags()
}

func (r *render) clearFlags() {
	r.flags = flags{}
}

func (r *render) parseFlags(string) error {
	return nil
}

func (r *render) render() error {
	// TODO(slongfield) Dispatch to different functions if the type is set.
	switch t := r.val.(type) {
	case string:
		r.buf.WriteString(r.val.(string))
		return nil
	case int:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int8:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int16:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int32:
		r.buf.WriteString(strconv.FormatInt(int64(t), 10))
		return nil
	case int64:
		r.buf.WriteString(strconv.FormatInt(t, 10))
		return nil
	case uint:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint8:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint16:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint32:
		r.buf.WriteString(strconv.FormatUint(uint64(t), 10))
		return nil
	case uint64:
		r.buf.WriteString(strconv.FormatUint(t, 10))
		return nil
	case reflect.Value:
		if t.IsValid() && t.CanInterface() {
			r.val = t.Interface()
			return r.render()
		}
		return r.renderValue(t)
	default:
		return fmt.Errorf("Unimplemented! %v %v", r.val, reflect.TypeOf(r.val).Kind())
	}
}

func (r *render) renderValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		return fmt.Errorf("Invalid value: %v", v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r.buf.WriteString(strconv.FormatInt(int64(v.Int()), 10))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r.buf.WriteString(strconv.FormatUint(uint64(v.Uint()), 10))
		return nil
	case reflect.String:
		r.buf.WriteString(v.String())
		return nil
	default:
		return fmt.Errorf("Unimplemented reflect type %v for %v ", v.Kind(), v)
	}
}

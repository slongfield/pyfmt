package pyfmt

import "errors"

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
	default:
		return errors.New("Unimplemented!")
	}
}

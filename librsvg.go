package librsvg

/*
#cgo pkg-config: librsvg-2.0
#include <cairo/cairo-pdf.h>
#include <librsvg/rsvg.h>
#include <stdint.h>
*/
import "C"
import (
	"unsafe"

	"github.com/kmhalpin/go-librsvg/pkg/glib"
	"github.com/ungerik/go-cairo"
)

type rsvg_RsvgHandle *C.RsvgHandle

type Handle struct {
	rsvgHandle rsvg_RsvgHandle
}

type Rectangle struct {
	x, y, width, height float64
}

func (h Handle) GetIntrinsicSizeInPixels() (width, height float64) {
	var rsvg_w, rsvg_h C.gdouble
	C.rsvg_handle_get_intrinsic_size_in_pixels(h.rsvgHandle, &rsvg_w, &rsvg_h)
	return float64(rsvg_w), float64(rsvg_h)
}

func (h Handle) RenderDocument(c cairo.Cairo_context, rectangle Rectangle) error {
	var rsvg_err *C.GError

	C.rsvg_handle_render_document(h.rsvgHandle, (*C.cairo_t)(unsafe.Pointer(c)), &C.RsvgRectangle{
		x:      C.gdouble(rectangle.x),
		y:      C.gdouble(rectangle.y),
		width:  C.gdouble(rectangle.width),
		height: C.gdouble(rectangle.height),
	}, &rsvg_err)
	if rsvg_err != nil {
		return glib.NewGError(uint32(rsvg_err.domain), int(rsvg_err.code), C.GoString(rsvg_err.message))
	}

	return nil
}

func NewHandle() *Handle {
	return &Handle{
		rsvgHandle: C.rsvg_handle_new(),
	}
}

func NewHandleFromData(data []byte) (*Handle, error) {
	var rsvg_err *C.GError

	rsvg_data := C.CBytes(data)
	defer C.free(rsvg_data)

	h := &Handle{
		rsvgHandle: C.rsvg_handle_new_from_data((*C.uchar)(rsvg_data), C.size_t(len(data)), &rsvg_err),
	}
	if rsvg_err != nil {
		return nil, glib.NewGError(uint32(rsvg_err.domain), int(rsvg_err.code), C.GoString(rsvg_err.message))
	}

	return h, nil
}

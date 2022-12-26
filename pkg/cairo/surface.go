package cairo

/*
#cgo pkg-config: cairo
#include <cairo/cairo-pdf.h>
cairo_status_t _go_cairo_write_func(void *closure, unsigned char* data, unsigned int length);
void _go_cairo_destroy_func(void *data);
*/
import "C"

import (
	"io"
	"runtime"
	"runtime/cgo"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/cairo"
)

// cairo_write_func_t
//
//export _go_cairo_write_func
func _go_cairo_write_func(closure *C.void, data *C.uchar, length C.uint) C.cairo_status_t {
	if writer, ok := cgo.Handle(*(*uintptr)(unsafe.Pointer(closure))).Value().(io.Writer); ok {
		if _, err := writer.Write(C.GoBytes(unsafe.Pointer(data), C.int(length))); err == nil {
			return C.CAIRO_STATUS_SUCCESS
		}
	}
	return C.CAIRO_STATUS_WRITE_ERROR
}

// cairo_destroy_func_t
//
//export _go_cairo_destroy_func
func _go_cairo_destroy_func(data *C.void) {
	cgo.Handle(*(*uintptr)(unsafe.Pointer(data))).Delete()
}

var (
	user_data_key_write_func = C.cairo_user_data_key_t{0}
)

func NewPDFSurfaceForStream(writer io.Writer, widthInPoints, heightInPoints float64) (*cairo.Surface, error) {
	h := cgo.NewHandle(writer)
	wptr := uintptr(h)

	s := C.cairo_pdf_surface_create_for_stream(
		C.cairo_write_func_t(C._go_cairo_write_func),
		unsafe.Pointer(&wptr), C.double(widthInPoints), C.double(heightInPoints))
	ws := cairo.WrapSurface(uintptr(unsafe.Pointer(s)))

	if status := ws.Status(); status != cairo.StatusSuccess {
		h.Delete()
		return nil, status
	}

	if status := cairo.Status(C.cairo_surface_set_user_data(s, &user_data_key_write_func,
		unsafe.Pointer(&wptr), C.cairo_destroy_func_t(C._go_cairo_destroy_func))); status != cairo.StatusSuccess {
		h.Delete()
		return nil, status
	}

	runtime.SetFinalizer(ws, (*cairo.Surface).Close)
	return ws, nil
}

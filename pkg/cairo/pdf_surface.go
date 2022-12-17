package cairo

/*
#cgo pkg-config: cairo
#include <cairo/cairo-pdf.h>
#include <stdint.h>

typedef unsigned char const* cairo_write_func_data_t;
extern cairo_status_t go_cairo_write_func(void *closure, cairo_write_func_data_t data, unsigned int length);
*/
import "C"

import (
	"io"
	"runtime/cgo"
	"unsafe"

	cairo "github.com/ungerik/go-cairo"
)

// cairo_write_func_t
//
//export go_cairo_write_func
func go_cairo_write_func(closure *C.void, data C.cairo_write_func_data_t, length C.uint) C.cairo_status_t {
	if writer, ok := cgo.Handle(*(*uintptr)(unsafe.Pointer(closure))).Value().(io.Writer); ok {
		if _, err := writer.Write(C.GoBytes(unsafe.Pointer(data), C.int(length))); err == nil {
			return C.CAIRO_STATUS_SUCCESS
		}
	}
	return C.CAIRO_STATUS_WRITE_ERROR
}

func NewSurfaceForStream(writer io.Writer, widthInPoints, heightInPoints float64) *cairo.Surface {
	wptr := uintptr(cgo.NewHandle(writer))
	s := C.cairo_pdf_surface_create_for_stream(
		C.cairo_write_func_t(C.go_cairo_write_func),
		unsafe.Pointer(&wptr), C.double(widthInPoints), C.double(heightInPoints))
	c := C.cairo_create(s)
	return cairo.NewSurfaceFromC(cairo.Cairo_surface(unsafe.Pointer(s)), cairo.Cairo_context(unsafe.Pointer(c)))
}

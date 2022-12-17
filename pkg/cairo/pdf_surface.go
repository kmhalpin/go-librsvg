package cairo

/*
#cgo pkg-config: librsvg-2.0
#include <cairo/cairo-pdf.h>
#include <librsvg/rsvg.h>
#include <stdint.h>
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"

	cairo "github.com/ungerik/go-cairo"
)

type Write func(closure interface{}, data []byte, length uint) cairo.Status

func NewSurfaceForStream(writeFunc Write, closure interface{}, widthInPoints, heightInPoints float64) *cairo.Surface {
	cl := cgo.NewHandle(closure)
	wr := cgo.NewHandle(writeFunc)
	s := C.cairo_pdf_surface_create_for_stream(
		C.cairo_write_func_t(unsafe.Pointer(&wr)),
		unsafe.Pointer(&cl), C.gdouble(widthInPoints), C.gdouble(heightInPoints))
	c := C.cairo_create(s)
	return cairo.NewSurfaceFromC(cairo.Cairo_surface(unsafe.Pointer(s)), cairo.Cairo_context(unsafe.Pointer(c)))
}

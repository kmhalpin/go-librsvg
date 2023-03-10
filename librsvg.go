package librsvg

/*
#cgo pkg-config: librsvg-2.0
#include <librsvg/rsvg.h>
*/
import "C"
import (
	"io"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/core/gerror"
	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

type rsvg_RsvgHandle *C.RsvgHandle

type Handle struct {
	rsvgHandle rsvg_RsvgHandle
}

type Rectangle struct {
	X, Y, Width, Height float64
}

func (h Handle) GetIntrinsicSizeInPixels() (width, height float64) {
	var rsvg_w, rsvg_h C.gdouble
	C.rsvg_handle_get_intrinsic_size_in_pixels(h.rsvgHandle, &rsvg_w, &rsvg_h)
	return float64(rsvg_w), float64(rsvg_h)
}

func (h Handle) RenderDocument(cairo *cairo.Context, rectangle Rectangle) error {
	var rsvg_err *C.GError
	if ok := C.rsvg_handle_render_document(h.rsvgHandle, *(**C.cairo_t)(unsafe.Pointer(cairo)), &C.RsvgRectangle{
		x:      C.gdouble(rectangle.X),
		y:      C.gdouble(rectangle.Y),
		width:  C.gdouble(rectangle.Width),
		height: C.gdouble(rectangle.Height),
	}, &rsvg_err); ok == 0 && rsvg_err != nil {
		return gerror.Take(unsafe.Pointer(rsvg_err))
	}

	return nil
}

func (h Handle) RenderElement(cairo *cairo.Context, id string, rectangle Rectangle) error {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	var rsvg_err *C.GError
	if ok := C.rsvg_handle_render_element(h.rsvgHandle, *(**C.cairo_t)(unsafe.Pointer(cairo)),
		cid,
		&C.RsvgRectangle{
			x:      C.gdouble(rectangle.X),
			y:      C.gdouble(rectangle.Y),
			width:  C.gdouble(rectangle.Width),
			height: C.gdouble(rectangle.Height),
		}, &rsvg_err); ok == 0 && rsvg_err != nil {
		return gerror.Take(unsafe.Pointer(rsvg_err))
	}

	return nil
}

func (h Handle) RenderLayer(cairo *cairo.Context, id string, rectangle Rectangle) error {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	var rsvg_err *C.GError
	if ok := C.rsvg_handle_render_layer(h.rsvgHandle, *(**C.cairo_t)(unsafe.Pointer(cairo)),
		cid,
		&C.RsvgRectangle{
			x:      C.gdouble(rectangle.X),
			y:      C.gdouble(rectangle.Y),
			width:  C.gdouble(rectangle.Width),
			height: C.gdouble(rectangle.Height),
		}, &rsvg_err); ok == 0 && rsvg_err != nil {
		return gerror.Take(unsafe.Pointer(rsvg_err))
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
		return nil, gerror.Take(unsafe.Pointer(rsvg_err))
	}

	return h, nil
}

type HandleFlags int

const (
	// No flags are set.
	RSVG_HANDLE_FLAGS_NONE HandleFlags = C.RSVG_HANDLE_FLAGS_NONE
	// Disable safety limits in the XML parser. Libxml2 has several limits designed to keep malicious
	// XML content from consuming too much memory while parsing. For security reasons, this should
	// only be used for trusted input! Since: 2.40.3
	RSVG_HANDLE_FLAG_UNLIMITED HandleFlags = C.RSVG_HANDLE_FLAG_UNLIMITED
	// Use this if the Cairo surface to which you are rendering is a PDF, PostScript, SVG, or Win32
	// Printing surface. This will make librsvg and Cairo use the original, compressed data for
	// images in the final output, instead of passing uncompressed images. For example, this
	// will make the a resulting PDF file smaller and faster. Please see the Cairo documentation
	// for details.
	RSVG_HANDLE_FLAG_KEEP_IMAGE_DATA HandleFlags = C.RSVG_HANDLE_FLAG_KEEP_IMAGE_DATA
	// Source: https://gnome.pages.gitlab.gnome.org/librsvg/Rsvg-2.0/flags.HandleFlags.html
)

func NewHandleFromStreamSync(inputStream io.Reader, file gio.Filer, handleFlags HandleFlags, cancellable *gio.Cancellable) (*Handle, error) {
	var rsvg_err *C.GError

	in := gioutil.NewInputStream(inputStream)
	inptr := in.Native()

	var f *C.GFile = nil
	if file != nil {
		fptr := glib.BaseObject(file).Native()
		f = *(**C.GFile)(unsafe.Pointer(&fptr))
	}

	var cn *C.GCancellable = nil
	if cancellable != nil {
		cnptr := cancellable.Native()
		cn = *(**C.GCancellable)(unsafe.Pointer(&cnptr))
	}

	h := &Handle{
		rsvgHandle: C.rsvg_handle_new_from_stream_sync(*(**C.GInputStream)(unsafe.Pointer(&inptr)),
			f, C.RsvgHandleFlags(handleFlags), cn, &rsvg_err),
	}
	if rsvg_err != nil {
		return nil, gerror.Take(unsafe.Pointer(rsvg_err))
	}

	return h, nil
}

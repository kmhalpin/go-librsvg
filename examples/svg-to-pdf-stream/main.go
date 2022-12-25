package main

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/kmhalpin/go-librsvg"
	cairoExt "github.com/kmhalpin/go-librsvg/pkg/cairo"
)

func main() {
	open, err := os.Open("test.svg")
	if err != nil {
		panic(err)
	}
	defer open.Close()

	file, err := os.OpenFile("test.pdf", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	h, err := librsvg.NewHandleFromStreamSync(open, nil, librsvg.RSVG_HANDLE_FLAG_KEEP_IMAGE_DATA, nil)
	if err != nil {
		panic(err)
	}
	width, height := h.GetIntrinsicSizeInPixels()

	s, err := cairoExt.NewPDFSurfaceForStream(file, width, height)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	c := cairo.Create(s)
	defer c.Close()

	if err := h.RenderDocument(c, librsvg.Rectangle{
		X: 0, Y: 0, Width: width, Height: height,
	}); err != nil {
		panic(err)
	}
}

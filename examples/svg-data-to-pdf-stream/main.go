package main

import (
	"os"

	"github.com/kmhalpin/go-librsvg"
	"github.com/kmhalpin/go-librsvg/pkg/cairo"
)

func main() {
	data, err := os.ReadFile("test.svg")
	if err != nil {
		panic(err)
	}

	h, err := librsvg.NewHandleFromData(data)
	if err != nil {
		panic(err)
	}
	width, height := h.GetIntrinsicSizeInPixels()

	file, err := os.OpenFile("test.pdf", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	s := cairo.NewSurfaceForStream(file, width, height)

	h.RenderDocument(s, librsvg.Rectangle{
		X: 0, Y: 0, Width: width, Height: height,
	})
	s.Destroy()
}

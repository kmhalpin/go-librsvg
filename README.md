## go-librsvg

Go binding for the librsvg library, which allows you to render SVG images in Go. This package depends to [gotk4](https://github.com/diamondburned/gotk4/).

### Installation
Install go-librsvg:

    go get github.com/kmhalpin/go-librsvg

go-librsvg is a cgo package and depends to librsvg for binding. You need `gcc` and `librsvg-dev` to build your app. however to run the app it only requires `librsvg`.

This package also extends feature in [gotk4/pkg/cairo](https://github.com/diamondburned/gotk4/) to construct cairo PDF surface for stream.

    go get github.com/kmhalpin/go-librsvg/pkg/cairo

### Features
|Feature          |Status|
|-----------------|------|
|SVG From Data (memory)|✅|
|SVG From GIO Stream|✅|
|SVG From File|❌|

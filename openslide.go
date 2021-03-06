package openslide

/*
#cgo CFLAGS: -g -Wall -I${SRCDIR}/include
#cgo LDFLAGS: -L. -lopenslide
#include "openslide.h"
#include "openslide-features.h"
*/
import "C"

import (
	"bytes"
	"image"
	"io"
	"unsafe"
)

func openOpenSlide(path string) (*C.openslide_t, error) {
	//openslide_t * osr
	pathToSVS := C.CString(path)
	return C.openslide_open(pathToSVS), nil
}

func openslide_get_level_count(osr *C.openslide_t) C.int {
	return C.openslide_get_level_count(osr)
}

func openslide_get_level_dimensions(osr *C.openslide_t, level int) (int64, int64, error) {
	var w, h C.int64_t
	levelC := C.int(level)
	C.openslide_get_level_dimensions(osr, levelC, &w, &h)
	return int64(w), int64(h), nil
}

func openslide_read_region(osr *C.openslide_t, level int, x int64, y int64, w int64, h int64) ([]byte, error) {
	buf := make([]byte, w*h*4)
	levelC := C.int(level)
	xC := C.int64_t(x)
	yC := C.int64_t(y)
	wC := C.int64_t(w)
	hC := C.int64_t(h)
	C.openslide_read_region(osr, (*C.uint32_t)(unsafe.Pointer(&buf[0])), xC, yC, levelC, wC, hC)
	return buf, nil
}

//-------------------------------------------------
// Somer helper functions to make life easier

func NewRGBAImageFromRegionInLevel(osr *C.openslide_t, level int, x, y, width, height int64) (*image.RGBA, error) {
	var err error
	imgBuf, _ := openslide_read_region(osr, level, x, y, width, height)
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	img.Pix, err = openslidePreMultipliedARGB2PixUINT8RGBA(imgBuf)
	return img, err
}

func openslidePreMultipliedARGB2PixUINT8RGBA(imgBuff []byte) ([]uint8, error) {
	pix := []uint8{}
	var err error
	buff := bytes.NewBuffer(imgBuff)
	for err == nil {
		pixel := buff.Next(4)
		if len(pixel) < 4 {
			err = io.EOF
			break
		}
		//                R         G         B         A
		a := pixel[3]
		r := pixel[2]
		g := pixel[1]
		b := pixel[0]
		if a != 0 && a != 255 {
			r = r * 255 / a
			g = g * 255 / a
			b = b * 255 / a
		}
		pix = append(pix, r, g, b, a)
	}
	if err == io.EOF {
		err = nil
	}
	return pix, err
}

// Pixel struct example
// type Pixel struct {
// 	R int
// 	G int
// 	B int
// 	A int
// }

// // img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
// func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
// 	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
// }

// type argb2rgba struct {
// 	buff bytes.Buffer
// }

// func (a *argb2rgba) Copy2PixUINT8() ([]uint8, error) {
// 	pix := []uint8{}
// 	var err error
// 	for err == nil {
// 		pixel := a.buff.Next(4)
// 		if len(pixel) < 4 {
// 			err = io.EOF
// 			break
// 		}
// 		//                R         G         B         A
// 		a := pixel[3]
// 		r := pixel[2]
// 		g := pixel[1]
// 		b := pixel[0]
// 		if a != 0 && a != 255 {
// 			r = r * 255 / a
// 			g = g * 255 / a
// 			b = b * 255 / a
// 		}
// 		pix = append(pix, r, g, b, a)
// 	}
// 	if err == io.EOF {
// 		err = nil
// 	}
// 	return pix, err
// }

// // Get the bi-dimensional pixel array
// func getPixels(file io.Reader) ([][]Pixel, error) {
// 	img, _, err := image.Decode(file)

// 	if err != nil {
// 		return nil, err
// 	}

// 	bounds := img.Bounds()
// 	width, height := bounds.Max.X, bounds.Max.Y

// 	var pixels [][]Pixel
// 	for y := 0; y < height; y++ {
// 		var row []Pixel
// 		for x := 0; x < width; x++ {
// 			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
// 		}
// 		pixels = append(pixels, row)
// 	}

// 	return pixels, nil
// }

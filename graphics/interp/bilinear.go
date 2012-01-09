// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

import (
	"image"
	"math"
)

// TODO(crawshaw): When we have inlining, consider func (p *RGBA) Off(x, y) int
func offRGBA(src *image.RGBA, x, y float64) int {
	return (int(y)-src.Rect.Min.Y)*src.Stride + (int(x)-src.Rect.Min.X)*4
}

// BinterpRGBA bilinear interpolates (sx, sy) and writes it to (x, y).
func BinterpRGBA(dst *image.RGBA, x, y int, src *image.RGBA, sx, sy float64) {
	srcb := src.Bounds()
	swidth := float64(srcb.Dx())
	sheight := float64(srcb.Dy())

	minX := float64(srcb.Min.X)
	minY := float64(srcb.Min.Y)
	lowX := math.Floor(sx - 0.5)
	lowY := math.Floor(sy - 0.5)
	if lowX < minX {
		lowX = minX
	}
	if lowY < minY {
		lowY = minY
	}

	highX := math.Ceil(sx - 0.5)
	highY := math.Ceil(sy - 0.5)
	if highX >= swidth {
		highX = swidth - 1
	}
	if highY >= sheight {
		highY = sheight - 1
	}

	// In the variables below, the 0 suffix indicates top/left, and the
	// 1 suffix indicates bottom/right.

	// Array offsets for the surrounding pixels.
	off00 := offRGBA(src, lowX, lowY)
	off01 := offRGBA(src, highX, lowY)
	off10 := offRGBA(src, lowX, highY)
	off11 := offRGBA(src, highX, highY)

	// Center of each surrounding pixel.
	x00 := lowX + 0.5
	y00 := lowY + 0.5
	x01 := highX + 0.5
	y01 := lowY + 0.5
	x10 := lowX + 0.5
	y10 := highY + 0.5
	x11 := highX + 0.5
	y11 := highY + 0.5

	frac00 := 0.0
	frac01 := 0.0
	frac10 := 0.0
	frac11 := 0.0

	// Literally, edge cases. If we are close enough to the edge of
	// the image, curtail the interpolation sources.
	if lowX == highX && lowY == highY {
		frac00 = 1.0
	} else if sy <= 0.5 && sx <= 0.5 {
		frac00 = 1.0
	} else if sheight-sy <= 0.5 && swidth-sx <= 0.5 {
		frac11 = 1.0
	} else if sy <= 0.5 || lowY == highY {
		frac00 = x01 - sx
		frac01 = sx - x00
	} else if sx <= 0.5 || lowX == highX {
		frac00 = y10 - sy
		frac10 = sy - y00
	} else if sheight-sy <= 0.5 {
		frac10 = x11 - sx
		frac11 = sx - x10
	} else if swidth-sx <= 0.5 {
		frac01 = y11 - sy
		frac11 = sy - y01
	} else {
		frac00 = (x01 - sx) * (y10 - sy)
		frac01 = (sx - x00) * (y11 - sy)
		frac10 = (x11 - sx) * (sy - y00)
		frac11 = (sx - x10) * (sy - y01)
	}

	r := 0.0
	g := 0.0
	b := 0.0
	a := 0.0

	r += float64(src.Pix[off00+0]) * frac00
	g += float64(src.Pix[off00+1]) * frac00
	b += float64(src.Pix[off00+2]) * frac00
	a += float64(src.Pix[off00+3]) * frac00

	r += float64(src.Pix[off01+0]) * frac01
	g += float64(src.Pix[off01+1]) * frac01
	b += float64(src.Pix[off01+2]) * frac01
	a += float64(src.Pix[off01+3]) * frac01

	r += float64(src.Pix[off10+0]) * frac10
	g += float64(src.Pix[off10+1]) * frac10
	b += float64(src.Pix[off10+2]) * frac10
	a += float64(src.Pix[off10+3]) * frac10

	r += float64(src.Pix[off11+0]) * frac11
	g += float64(src.Pix[off11+1]) * frac11
	b += float64(src.Pix[off11+2]) * frac11
	a += float64(src.Pix[off11+3]) * frac11

	off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
	dst.Pix[off+0] = uint8(r + 0.5)
	dst.Pix[off+1] = uint8(g + 0.5)
	dst.Pix[off+2] = uint8(b + 0.5)
	dst.Pix[off+3] = uint8(a + 0.5)
}

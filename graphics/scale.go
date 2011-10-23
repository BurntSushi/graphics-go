// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
)

func scaleRGBA(src, dst *image.RGBA) {
	srcb := src.Bounds()
	b := dst.Bounds()
	xratio := float64(srcb.Dx()) / float64(b.Dx())
	yratio := float64(srcb.Dy()) / float64(b.Dy())

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// The center of the destination pixel we are interpolating,
			// projected onto the source image.
			sx := (float64(x) + 0.5) * xratio
			sy := (float64(y) + 0.5) * yratio
			binterpRGBA(dst, x, y, src, sx, sy)
		}
	}
}

// Scale produces a scaled version of the image using bilinear interpolation.
func Scale(dst draw.Image, src image.Image) {
	if dst == nil || src == nil {
		return
	}

	b := dst.Bounds()
	if src.Bounds().Empty() || b.Empty() {
		return
	}

	srcRgba, ok := src.(*image.RGBA)
	if !ok {
		b := src.Bounds()
		srcRgba = image.NewRGBA(b)
		draw.Draw(srcRgba, b, src, b.Min, draw.Src)
	}

	dstRgba, ok := dst.(*image.RGBA)
	if !ok {
		dstRgba = image.NewRGBA(b)
	}

	scaleRGBA(srcRgba, dstRgba)

	if !ok {
		draw.Draw(dst, b, dstRgba, b.Min, draw.Src)
	}
}

// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
	"math"
)

func rotateRGBA(dst, src *image.RGBA, angle float64) {
	dstb := dst.Bounds()
	srcb := src.Bounds()
	dx := float64(dstb.Dx()) / 2
	dy := float64(dstb.Dy()) / 2
	sx := float64(srcb.Dx()) / 2
	sy := float64(srcb.Dy()) / 2
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// line up the center of dst with the center of src.
	// for each pixel of dst, rotate onto src and linear interpolate
	for y := dstb.Min.Y; y < dstb.Max.Y; y++ {
		for x := dstb.Min.X; x < dstb.Max.X; x++ {
			// Convert position to be zeroed on the center of dst.
			ldx := (float64(x) + 0.5 - dx)
			ldy := (float64(y) + 0.5 - dy)

			// Rotate position on dst to a position on src.
			// (Using the standard cartesian plane rotation matrix.)
			lx := ldx*cos + ldy*sin
			ly := ldy*cos - ldx*sin

			// Convert back to pixel position on src.
			px := lx + sx
			py := ly + sy

			if px < float64(srcb.Min.X) || px >= float64(srcb.Max.X) {
				continue
			} else if py < float64(srcb.Min.Y) || py >= float64(srcb.Max.Y) {
				continue
			}

			// Determine Color by Linear Interpolation.
			binterpRGBA(dst, x, y, src, px, py)
		}
	}
}

// RotateOptions are the rotation parameters.
// Angle is the angle, in radians, to rotate the image clockwise.
type RotateOptions struct {
	Angle float64
}

// Rotate produces a rotated version of src, drawn onto dst.
func Rotate(dst draw.Image, src image.Image, opt *RotateOptions) {
	if dst == nil || src == nil {
		return
	}

	angle := 0.0
	if opt != nil {
		angle = opt.Angle
	}

	srcRgba, ok := src.(*image.RGBA)
	if !ok {
		b := src.Bounds()
		srcRgba = image.NewRGBA(b)
		draw.Draw(srcRgba, b, src, b.Min, draw.Src)
	}

	b := dst.Bounds()
	dstRgba, ok := dst.(*image.RGBA)
	if !ok {
		dstRgba = image.NewRGBA(b)
	}

	rotateRGBA(dstRgba, srcRgba, angle)

	if !ok {
		draw.Draw(dst, b, dstRgba, b.Min, draw.Src)
	}
}

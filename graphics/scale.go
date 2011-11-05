// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
	"os"
)

// Scale produces a scaled version of the image using bilinear interpolation.
func Scale(dst draw.Image, src image.Image) os.Error {
	if dst == nil {
		return os.NewError("graphics: dst is nil")
	}
	if src == nil {
		return os.NewError("graphics: src is nil")
	}

	b := dst.Bounds()
	srcb := src.Bounds()
	if b.Empty() || srcb.Empty() {
		return nil
	}
	sx := float64(b.Dx()) / float64(srcb.Dx())
	sy := float64(b.Dy()) / float64(srcb.Dy())
	return I.Scale(sx, sy).Transform(dst, src)
}

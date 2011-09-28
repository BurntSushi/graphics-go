// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
	"math"
)

func blurRGBA(dst, src *image.RGBA, sd float64, size int) {
	// kernel is the Gaussian kernel weights, scaled so that kernel[0] == 1<<16.
	kernel := make([]uint64, size+1)
	for i := 0; i <= size; i++ {
		x := float64(i) / sd
		kernel[i] = uint64(math.Pow(1/math.SqrtE, x*x) * (1 << 16))
	}
	// ksum is the 2-D kernel sum, which is the square of the 1-D kernel sum.
	var kSum uint64
	for i := 1; i <= size; i++ {
		kSum += kernel[i]
	}
	kSum = 2*kSum + kernel[0]
	kSum *= kSum

	// buf holds the result of vertically blurring src.
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	buf := make([]uint64, width*height*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a uint64
			// k0 is the kernel weight for the center pixel. This may be greater
			// than kernel[0], near the boundary of the source image, to avoid
			// vignetting. off is the src.Pix offset of the center pixel.
			k0, off := kernel[0], (y-bounds.Min.Y)*src.Stride+(x-bounds.Min.X)*4

			// Add the pixels from above.
			for i := 1; i <= size; i++ {
				k := kernel[i]
				if y-i < bounds.Min.Y {
					k0 += k
				} else {
					o := off - i*src.Stride
					r += uint64(src.Pix[o+0]) * k
					g += uint64(src.Pix[o+1]) * k
					b += uint64(src.Pix[o+2]) * k
					a += uint64(src.Pix[o+3]) * k
				}
			}

			// Add the pixels from below.
			for i := 1; i <= size; i++ {
				k := kernel[i]
				if y+i >= bounds.Max.Y {
					k0 += k
				} else {
					o := off + i*src.Stride
					r += uint64(src.Pix[o+0]) * k
					g += uint64(src.Pix[o+1]) * k
					b += uint64(src.Pix[o+2]) * k
					a += uint64(src.Pix[o+3]) * k
				}
			}

			// Add the central pixel.
			r += uint64(src.Pix[off+0]) * k0
			g += uint64(src.Pix[off+1]) * k0
			b += uint64(src.Pix[off+2]) * k0
			a += uint64(src.Pix[off+3]) * k0

			// Write to buf.
			o := (y-bounds.Min.Y)*width*4 + (x-bounds.Min.X)*4
			buf[o+0] = r
			buf[o+1] = g
			buf[o+2] = b
			buf[o+3] = a
		}
	}

	// dst holds the result of horizontally blurring buf.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a uint64
			k0, off := kernel[0], y*width*4+x*4

			// Add the pixels from the left.
			for i := 1; i <= size; i++ {
				k := kernel[i]
				if x-i < 0 {
					k0 += k
				} else {
					o := off - i*4
					r += buf[o+0] * k
					g += buf[o+1] * k
					b += buf[o+2] * k
					a += buf[o+3] * k
				}
			}

			// Add the pixels from the right.
			for i := 1; i <= size; i++ {
				k := kernel[i]
				if x+i >= width {
					k0 += k
				} else {
					o := off + i*4
					r += buf[o+0] * k
					g += buf[o+1] * k
					b += buf[o+2] * k
					a += buf[o+3] * k
				}
			}

			// Add the central pixel.
			r += buf[off+0] * k0
			g += buf[off+1] * k0
			b += buf[off+2] * k0
			a += buf[off+3] * k0

			// Write to dst, scaling back to the range [0, 255].
			dstOff := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
			dst.Pix[dstOff+0] = uint8((r + kSum/2) / kSum)
			dst.Pix[dstOff+1] = uint8((g + kSum/2) / kSum)
			dst.Pix[dstOff+2] = uint8((b + kSum/2) / kSum)
			dst.Pix[dstOff+3] = uint8((a + kSum/2) / kSum)
		}
	}
}

// DefaultStdDev is the default blurring parameter.
var DefaultStdDev = 0.5

// BlurOptions are the blurring parameters.
// StdDev is the standard deviation of the normal, higher is blurrier.
// Size is the size of the kernel. If zero, it is set to Ceil(6 * StdDev).
type BlurOptions struct {
	StdDev float64
	Size   int
}

// Blur produces a blurred version of the image, using a Gaussian blur.
func Blur(dst draw.Image, src image.Image, opt *BlurOptions) {
	sd := DefaultStdDev
	size := 0

	if opt != nil {
		sd = opt.StdDev
		size = opt.Size
	}

	if size == 0 {
		size = int(math.Ceil(sd * 6))
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

	blurRGBA(dstRgba, srcRgba, sd, size)

	if !ok {
		draw.Draw(dst, b, dstRgba, b.Min, draw.Src)
	}
}

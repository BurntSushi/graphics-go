// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
)

// clamp clamps x to the range [x0, x1].
func clamp(x, x0, x1 float64) float64 {
	if x < x0 {
		return x0
	}
	if x > x1 {
		return x1
	}
	return x
}

// Kernel is a square matrix that defines a convolution.
type Kernel interface {
	// TODO: document how the slice is laid out.
	Weights() []float64
}

// SeparableKernel is a linearly separable, square convolution kernel.
// X and Y are the per-axis weights. Each slice must be the same length, and
// have an odd length. The middle element of each slice is the weight for the
// central pixel. For example, the horizontal Sobel kernel is:
//	sobelX := &SeparableKernel{
//		X: []float64{-1, 0, +1},
//		Y: []float64{1, 2, 1},
//	}
type SeparableKernel struct {
	X, Y []float64
}

func (k *SeparableKernel) Weights() []float64 {
	panic("unimplemented")
}

func convolveRGBASep(dst *image.RGBA, src *image.RGBA, k *SeparableKernel) {
	if len(k.X) != len(k.Y) || len(k.X)%2 != 1 {
		// TODO: return an error.
		panic("invalid kernel")
	}
	radius := (len(k.X) - 1) / 2

	// buf holds the result of vertically blurring src.
	bounds := dst.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	buf := make([]float64, width*height*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64
			// k0 is the kernel weight for the center pixel. This may be greater
			// than kernel[0], near the boundary of the source image, to avoid
			// vignetting. off is the src.Pix offset of the center pixel.
			k0, off := k.X[radius], (y-bounds.Min.Y)*src.Stride+(x-bounds.Min.X)*4

			// Add the pixels from above.
			for i := 1; i <= radius; i++ {
				f := k.Y[radius-i]
				if y-i < bounds.Min.Y {
					k0 += f
				} else {
					o := off - i*src.Stride
					r += float64(src.Pix[o+0]) * f
					g += float64(src.Pix[o+1]) * f
					b += float64(src.Pix[o+2]) * f
					a += float64(src.Pix[o+3]) * f
				}
			}

			// Add the pixels from below.
			for i := 1; i <= radius; i++ {
				f := k.Y[radius+i]
				if y+i >= bounds.Max.Y {
					k0 += f
				} else {
					o := off + i*src.Stride
					r += float64(src.Pix[o+0]) * f
					g += float64(src.Pix[o+1]) * f
					b += float64(src.Pix[o+2]) * f
					a += float64(src.Pix[o+3]) * f
				}
			}

			// Add the central pixel.
			r += float64(src.Pix[off+0]) * k0
			g += float64(src.Pix[off+1]) * k0
			b += float64(src.Pix[off+2]) * k0
			a += float64(src.Pix[off+3]) * k0

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
			var r, g, b, a float64
			k0, off := k.X[radius], y*width*4+x*4

			// Add the pixels from the left.
			for i := 1; i <= radius; i++ {
				f := k.X[radius-i]
				if x-i < 0 {
					k0 += f
				} else {
					o := off - i*4
					r += buf[o+0] * f
					g += buf[o+1] * f
					b += buf[o+2] * f
					a += buf[o+3] * f
				}
			}

			// Add the pixels from the right.
			for i := 1; i <= radius; i++ {
				f := k.X[radius+i]
				if x+i >= width {
					k0 += f
				} else {
					o := off + i*4
					r += buf[o+0] * f
					g += buf[o+1] * f
					b += buf[o+2] * f
					a += buf[o+3] * f
				}
			}

			// Add the central pixel.
			r += buf[off+0] * k0
			g += buf[off+1] * k0
			b += buf[off+2] * k0
			a += buf[off+3] * k0

			// Write to dst, clamping to the range [0, 255].
			dstOff := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
			dst.Pix[dstOff+0] = uint8(clamp(r+0.5, 0, 255))
			dst.Pix[dstOff+1] = uint8(clamp(g+0.5, 0, 255))
			dst.Pix[dstOff+2] = uint8(clamp(b+0.5, 0, 255))
			dst.Pix[dstOff+3] = uint8(clamp(a+0.5, 0, 255))
		}
	}
}

// Convolve produces dst by applying the convolution kernel k to src.
func Convolve(dst draw.Image, src image.Image, k Kernel) {
	b := dst.Bounds()
	dstRgba, ok := dst.(*image.RGBA)
	if !ok {
		dstRgba = image.NewRGBA(b)
	}

	srcRgba, ok := src.(*image.RGBA)
	if !ok {
		b := src.Bounds()
		srcRgba = image.NewRGBA(b)
		draw.Draw(srcRgba, b, src, b.Min, draw.Src)
	}

	switch k := k.(type) {
	case *SeparableKernel:
		convolveRGBASep(dstRgba, srcRgba, k)
	default:
		panic("unimplemented")
	}

	if !ok {
		draw.Draw(dst, b, dstRgba, b.Min, draw.Src)
	}
}

// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"image/draw"
	"math"
)

func circularGauss(x float64, sd float64) float64 {
	sdSq := sd * sd
	// If we were calculating the standard 1D Gaussian, this would
	// be 1/math.Sqrt(...). However this is a single dimension of
	// the circularly symmetric Gaussian. So no Sqrt.
	coeff := 1 / (2 * math.Pi * sdSq)
	return coeff * math.Pow(math.E, -1*math.Pow(x, 2)/(2*sdSq))
}

func blurRGBA(dst, src *image.RGBA, sd float64, size int) {
	kernel := make([]float64, size+1)
	for i := 0; i <= size; i++ {
		kernel[i] = circularGauss(float64(i), sd)
	}

	kernelSum := 0.0
	for i := 1; i <= size; i++ {
		kernelSum += kernel[i]
	}
	kernelSum = 2*kernelSum + kernel[0]
	// The 2-D kernel sum is the square of the 1-D kernel sum.
	kernelSum *= kernelSum

	bounds := src.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	ds1 := make([]float64, width*height*4)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b, a float64
			off := (y-src.Rect.Min.Y)*src.Stride + (x-src.Rect.Min.X)*4

			// When we fail to copy non-existant data near the boundaries of the
			// image, maintain an adjustment factor and use it to take more from
			// the original pixel. Avoids vignetting.
			adjust := float64(0)

			// Add blur from above.
			for i := 1; i <= size; i++ {
				if y-i < bounds.Min.Y {
					adjust += kernel[i]
				} else {
					o := off - i*src.Stride
					factor := kernel[i]
					r += float64(src.Pix[o+0]) * factor
					g += float64(src.Pix[o+1]) * factor
					b += float64(src.Pix[o+2]) * factor
					a += float64(src.Pix[o+3]) * factor
				}
			}

			// Add blur from below.
			for i := 1; i <= size; i++ {
				if y+i >= bounds.Max.Y {
					adjust += kernel[i]
				} else {
					o := off + i*src.Stride
					factor := kernel[i]
					r += float64(src.Pix[o+0]) * factor
					g += float64(src.Pix[o+1]) * factor
					b += float64(src.Pix[o+2]) * factor
					a += float64(src.Pix[o+3]) * factor
				}
			}

			// Add a fraction of the original pixel.
			factor := kernel[0] + adjust
			r += float64(src.Pix[off+0]) * factor
			g += float64(src.Pix[off+1]) * factor
			b += float64(src.Pix[off+2]) * factor
			a += float64(src.Pix[off+3]) * factor

			o := (y-bounds.Min.Y)*width*4 + (x-bounds.Min.X)*4
			ds1[o+0] = r
			ds1[o+1] = g
			ds1[o+2] = b
			ds1[o+3] = a
		}
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a float64
			adjust := float64(0)
			off := y*width*4 + x*4

			// Add blur from the left.
			for i := 1; i <= size; i++ {
				if x-i < 0 {
					adjust += kernel[i]
				} else {
					o := off - i*4
					factor := kernel[i]
					r += ds1[o+0] * factor
					g += ds1[o+1] * factor
					b += ds1[o+2] * factor
					a += ds1[o+3] * factor
				}
			}

			// Add blur from the right.
			for i := 1; i <= size; i++ {
				if x+i >= width {
					adjust += kernel[i]
				} else {
					o := off + i*4
					factor := kernel[i]
					r += ds1[o+0] * factor
					g += ds1[o+1] * factor
					b += ds1[o+2] * factor
					a += ds1[o+3] * factor
				}
			}

			// Add a fraction of the original pixel.
			factor := kernel[0] + adjust
			r += ds1[off+0] * factor
			g += ds1[off+1] * factor
			b += ds1[off+2] * factor
			a += ds1[off+3] * factor

			// When the kernel does not sum to 1.0, the resulting image has
			// less intensity than the original. Dividing through by the sum
			// corrects this.
			r /= kernelSum
			g /= kernelSum
			b /= kernelSum
			a /= kernelSum

			dstOff := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
			dst.Pix[dstOff+0] = uint8(math.Fmin(255.0, r+0.5))
			dst.Pix[dstOff+1] = uint8(math.Fmin(255.0, g+0.5))
			dst.Pix[dstOff+2] = uint8(math.Fmin(255.0, b+0.5))
			dst.Pix[dstOff+3] = uint8(math.Fmin(255.0, a+0.5))
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

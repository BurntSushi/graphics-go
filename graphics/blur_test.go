// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics_test

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"

	"graphics-go.googlecode.com/hg/graphics"
)

func delta(u0, u1 uint32) int {
	d := int(u0) - int(u1)
	if d < 0 {
		return -d
	}
	return d
}

func withinTolerance(c0, c1 image.Color, tolerance int) bool {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r := delta(r0, r1)
	g := delta(g0, g1)
	b := delta(b0, b1)
	a := delta(a0, a1)
	return r <= tolerance && g <= tolerance && b <= tolerance && a <= tolerance
}

func sprintBox(box []byte, width, height int) string {
	buf := bytes.NewBuffer(nil)
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Fprintf(buf, " 0x%02x,", box[i])
			i++
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

type blurOneColorTest struct {
	desc   string
	width  int
	height int
	opt    *graphics.BlurOptions
	src    []uint8
	res    []uint8
}

var blurOneColorTests = []blurOneColorTest{
	{
		"1x1-blank", 1, 1,
		&graphics.BlurOptions{0.83, 1},
		[]uint8{0xff},
		[]uint8{0xff},
	},
	{
		"1x1-spreadblank", 1, 1,
		&graphics.BlurOptions{0.83, 2},
		[]uint8{0xff},
		[]uint8{0xff},
	},
	{
		"3x3-blank", 3, 3,
		&graphics.BlurOptions{0.83, 2},
		[]uint8{
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
		},
		[]uint8{
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
			0xff, 0xff, 0xff,
		},
	},
	{
		"3x3-dot", 3, 3,
		&graphics.BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00,
			0x00, 0xff, 0x00,
			0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x03, 0x00,
			0x03, 0xf2, 0x03,
			0x00, 0x03, 0x00,
		},
	},
	{
		"5x5-dot", 5, 5,
		&graphics.BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xff, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00,
			0x00, 0x03, 0xf2, 0x03, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
		"5x5-dot-spread", 5, 5,
		&graphics.BlurOptions{0.85, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xff, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x10, 0x20, 0x10, 0x00,
			0x00, 0x20, 0x40, 0x20, 0x00,
			0x00, 0x10, 0x20, 0x10, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
	{
		"4x4-box", 4, 4,
		&graphics.BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0xff, 0xff, 0x00,
			0x00, 0xff, 0xff, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x03, 0x03, 0x00,
			0x03, 0xf8, 0xf8, 0x03,
			0x03, 0xf8, 0xf8, 0x03,
			0x00, 0x03, 0x03, 0x00,
		},
	},
	{
		"5x5-twodots", 5, 5,
		&graphics.BlurOptions{0.34, 1},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x96, 0x00, 0x96, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
		[]uint8{
			0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x02, 0x00,
			0x02, 0x8e, 0x04, 0x8e, 0x02,
			0x00, 0x02, 0x00, 0x02, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00,
		},
	},
}

func TestBlurOneColor(t *testing.T) {
	for _, oc := range blurOneColorTests {
		img := image.NewRGBA(image.Rect(0, 0, oc.width, oc.height))
		i := 0
		for y := 0; y < oc.height; y++ {
			for x := 0; x < oc.width; x++ {
				img.SetRGBA(x, y, image.RGBAColor{R: oc.src[i]})
				i++
			}
		}

		imgBlur, ok := graphics.Blur(img, oc.opt).(*image.RGBA)
		if !ok {
			t.Errorf("%s: blurred image is not RGBA", oc.desc)
			continue
		}

		i = 0
		res := make([]byte, len(oc.res))
		for y := 0; y < oc.height; y++ {
			for x := 0; x < oc.width; x++ {
				off := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)*4
				res[i] = imgBlur.Pix[off]
				i++
			}
		}

		if !bytes.Equal(res, oc.res) {
			got := sprintBox(res, oc.width, oc.height)
			want := sprintBox(oc.res, oc.width, oc.height)
			t.Errorf("%s:\n got\n%s\n want\n%s", oc.desc, got, want)
			continue
		}
	}
}

func benchBlur(b *testing.B, bounds image.Rectangle) {
	b.StopTimer()

	// Construct a fuzzy image.
	img := image.NewRGBA(bounds)
	dy := bounds.Dy()
	dx := bounds.Dx()
	for y := 0; y < dy; y++ {
		for x := 0; x < dx/4; x++ {
			img.SetRGBA(x, y, image.RGBAColor{
				uint8(5 * x % 0x100),
				uint8(7 * y % 0x100),
				uint8((7*x + 5*y) % 0x100),
				0xff,
			})
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		graphics.Blur(img, &graphics.BlurOptions{0.84, 3})
	}
}

func loadImage(path string) (img image.Image, err os.Error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

func BenchmarkBlur400x400x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 400))
}

// Exactly twice the pixel count of 400x400.
func BenchmarkBlur400x800x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 800))
}

// Exactly twice the pixel count of 400x800
func BenchmarkBlur400x1600x3(b *testing.B) {
	benchBlur(b, image.Rect(0, 0, 400, 1600))
}

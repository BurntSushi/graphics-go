// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"
)

type transformOneColorTest struct {
	desc      string
	dstWidth  int
	dstHeight int
	srcWidth  int
	srcHeight int
	opt       interface{}
	src       []uint8
	res       []uint8
}

func (oc *transformOneColorTest) newSrc() *image.RGBA {
	b := image.Rect(0, 0, oc.srcWidth, oc.srcHeight)
	src := image.NewRGBA(b)
	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			src.SetRGBA(x, y, image.RGBAColor{
				R: oc.src[i],
				G: oc.src[i],
				B: oc.src[i],
				A: oc.src[i],
			})
			i++
		}
	}
	return src
}

func (oc *transformOneColorTest) newDst() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, oc.dstWidth, oc.dstHeight))
}

func checkTransformTest(t *testing.T, oc *transformOneColorTest, dst, src *image.RGBA) bool {
	for ch := 0; ch < 4; ch++ {
		i := 0
		res := make([]byte, len(oc.res))
		for y := 0; y < oc.dstHeight; y++ {
			for x := 0; x < oc.dstWidth; x++ {
				off := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
				res[i] = dst.Pix[off+ch]
				i++
			}
		}

		if !bytes.Equal(res, oc.res) {
			got := sprintBox(res, oc.dstWidth, oc.dstHeight)
			want := sprintBox(oc.res, oc.dstWidth, oc.dstHeight)
			t.Errorf("%s: ch=%d\n got\n%s\n want\n%s", oc.desc, ch, got, want)
			return false
		}
	}

	return true
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

func delta(u0, u1 uint32) int {
	d := int(u0) - int(u1)
	if d < 0 {
		return -d
	}
	return d
}

func withinTolerance(c0, c1 image.Color, tol int) bool {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	r := delta(r0, r1)
	g := delta(g0, g1)
	b := delta(b0, b1)
	a := delta(a0, a1)
	return r <= tol && g <= tol && b <= tol && a <= tol
}

func imageWithinTolerance(m0, m1 image.Image, tol int) os.Error {
	b0 := m0.Bounds()
	b1 := m1.Bounds()
	if !b0.Eq(b1) {
		return os.NewError(fmt.Sprintf("got bounds %v want %v", b0, b1))
	}

	for y := b0.Min.Y; y < b0.Max.Y; y++ {
		for x := b0.Min.X; x < b0.Max.X; x++ {
			c0 := m0.At(x, y)
			c1 := m1.At(x, y)
			if !withinTolerance(c0, c1, tol) {
				e := fmt.Sprintf("got %v want %v at (%d, %d)", c0, c1, x, y)
				return os.NewError(e)
			}
		}
	}
	return nil
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

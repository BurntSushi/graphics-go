// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"testing"
	"image"

	_ "image/png"
)

var scaleOneColorTests = []transformOneColorTest{
	{
		"down-half",
		1, 1,
		2, 2,
		nil,
		[]uint8{
			0x80, 0x00,
			0x00, 0x80,
		},
		[]uint8{
			0x40,
		},
	},
	{
		"up-double",
		4, 4,
		2, 2,
		nil,
		[]uint8{
			0x80, 0x00,
			0x00, 0x80,
		},
		[]uint8{
			0x80, 0x60, 0x20, 0x00,
			0x60, 0x50, 0x30, 0x20,
			0x20, 0x30, 0x50, 0x60,
			0x00, 0x20, 0x60, 0x80,
		},
	},
	{
		"up-doublewidth",
		4, 2,
		2, 2,
		nil,
		[]uint8{
			0x80, 0x00,
			0x00, 0x80,
		},
		[]uint8{
			0x80, 0x60, 0x20, 0x00,
			0x00, 0x20, 0x60, 0x80,
		},
	},
	{
		"up-doubleheight",
		2, 4,
		2, 2,
		nil,
		[]uint8{
			0x80, 0x00,
			0x00, 0x80,
		},
		[]uint8{
			0x80, 0x00,
			0x60, 0x20,
			0x20, 0x60,
			0x00, 0x80,
		},
	},
	{
		"up-partial",
		3, 3,
		2, 2,
		nil,
		[]uint8{
			0x80, 0x00,
			0x00, 0x80,
		},
		[]uint8{
			0x80, 0x40, 0x00,
			0x40, 0x40, 0x40,
			0x00, 0x40, 0x80,
		},
	},
}

func TestScaleOneColor(t *testing.T) {
	for _, oc := range scaleOneColorTests {
		dst := oc.newDst()
		src := oc.newSrc()
		Scale(dst, src)

		if !checkTransformTest(t, &oc, dst) {
			continue
		}
	}
}

func TestScaleNil(t *testing.T) {
	Scale(nil, nil)
}

func TestScaleEmpty(t *testing.T) {
	empty := image.NewRGBA(image.Rect(0, 0, 0, 0))
	Scale(empty, empty)
}

func TestScaleGopher(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 100, 150))

	src, err := loadImage("../testdata/gopher.png")
	if err != nil {
		t.Error(err)
		return
	}

	// Down-sample.
	Scale(dst, src)
	cmp, err := loadImage("../testdata/gopher-100x150.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = imageWithinTolerance(dst, cmp, 0)
	if err != nil {
		t.Error(err)
		return
	}

	// Up-sample.
	dst = image.NewRGBA(image.Rect(0, 0, 500, 750))
	Scale(dst, src)
	cmp, err = loadImage("../testdata/gopher-500x750.png")
	if err != nil {
		t.Error(err)
		return
	}
	err = imageWithinTolerance(dst, cmp, 0)
	if err != nil {
		t.Error(err)
		return
	}
}

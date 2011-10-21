// Copyright 2011 The Graphics-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

import (
	"image"
	"reflect"
	"testing"

	_ "image/png"
)

func TestSeparableWeights(t *testing.T) {
	sobelXFull := []float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}
	sobelXSep := &SeparableKernel{
		X: []float64{-1, 0, +1},
		Y: []float64{1, 2, 1},
	}
	w := sobelXSep.Weights()
	if !reflect.DeepEqual(w, sobelXFull) {
		t.Errorf("got %v want %v", w, sobelXFull)
	}
}

func TestConvolve(t *testing.T) {
	kernFull, err := NewKernel([]float64{
		0, 0, 0,
		1, 1, 1,
		0, 0, 0,
	})
	if err != nil {
		t.Fatal(err)
	}

	kernSep := &SeparableKernel{
		X: []float64{1, 1, 1},
		Y: []float64{0, 1, 0},
	}

	src, err := loadImage("../testdata/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	b := src.Bounds()

	sep := image.NewRGBA(b)
	Convolve(sep, src, kernSep)

	full := image.NewRGBA(b)
	Convolve(full, src, kernFull)

	err = imageWithinTolerance(sep, full, 0x101)
	if err != nil {
		t.Fatal(err)
	}
}

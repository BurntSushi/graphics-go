package graphics

import (
	"image"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/ieee0824/graphics-go/graphics/interp"
)

func init() {
	now := time.Now()
	rand.Seed(now.UnixNano())
}

func toRadian(n int) float64 {
	return float64(n) * math.Pi / 180.0
}

var affines map[int]*Affine = map[int]*Affine{
	1: I,
	2: I.Scale(-1, 1),
	3: I.Scale(-1, -1),
	4: I.Scale(1, -1),
	5: I.Rotate(toRadian(90)).Scale(-1, 1),
	6: I.Rotate(toRadian(90)),
	7: I.Rotate(toRadian(-90)).Scale(-1, 1),
	8: I.Rotate(toRadian(-90)),
}

func genRandamAffine() *Affine {
	return &Affine{
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	}
}

func BenchmarkAffine_Mul(b *testing.B) {
	b.StopTimer()
	aAffine := genRandamAffine()
	bAffine := genRandamAffine()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bAffine = aAffine.Mul(bAffine)
	}
}

func BenchmarkAffine_TransformCenter(b *testing.B) {
	b.StopTimer()
	r := image.Rect(0, 0, 1024, 1024)
	aImg := image.NewRGBA64(r)
	bImg := image.NewRGBA64(r)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%9 == 0 {
			continue
		}
		affines[i%9].TransformCenter(aImg, bImg, interp.Bilinear)
	}
}

func BenchmarkAffine_CenterFit(b *testing.B) {
	b.StopTimer()
	r := image.Rect(0, 0, 1024, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%9 == 0 {
			continue
		}
		affines[i%9].CenterFit(r, r)
	}
}

func BenchmarkAffine_Transform(b *testing.B) {
	b.StopTimer()
	r := image.Rect(0, 0, 1024, 1024)
	aImg := image.NewRGBA64(r)
	bImg := image.NewRGBA64(r)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%9 == 0 {
			continue
		}
		affines[i%9].Transform(aImg, bImg, interp.Bilinear)
	}

}

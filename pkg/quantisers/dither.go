package quantisers

import (
	"github.com/fiwippi/go-quantise/pkg/colours"
	"image"
	"image/color"
	"math"
)

type DitherType int

const (
	NoDither DitherType = iota
	FloydSteinberg
	Bayer4x4
	Bayer8x8
)

// No Dither
func noDitherSingle(cimg *image.RGBA, c color.Palette) *image.RGBA {
	bounds := cimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := cimg.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			greyscaleLevel := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			Y := c[0].(color.Gray).Y

			if greyscaleLevel <= Y {
				cimg.Set(x, y, BLACK)
			} else {
				cimg.Set(x, y, WHITE)
			}
		}
	}

	return cimg
}

func noDitherMulti(cimg *image.RGBA, c color.Palette) *image.RGBA {
	bounds := cimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			cimg.Set(x, y, c.Convert(cimg.At(x, y)))
		}
	}

	return cimg
}

// Floyd-steinberg dithering https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering
func floydSteinbergDither(cimg *image.RGBA, c color.Palette) *image.RGBA {
	bounds := cimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			oldColour := cimg.At(x, y)
			newColour := c.Convert(oldColour)
			cimg.Set(x, y, newColour)

			rErr, gErr, bErr := fsQuantisedErrors(newColour, oldColour)
			cimg.Set(x+1, y, fsDiffuseErrors(x+1, y, cimg, rErr, gErr, bErr, 7.0/16))
			cimg.Set(x-1, y+1, fsDiffuseErrors(x-1, y+1, cimg, rErr, gErr, bErr, 3.0/16))
			cimg.Set(x, y+1, fsDiffuseErrors(x, y+1, cimg, rErr, gErr, bErr, 5.0/16))
			cimg.Set(x+1, y+1, fsDiffuseErrors(x+1, y+1, cimg, rErr, gErr, bErr, 1.0/16))
		}
	}

	return cimg
}

func fsDiffuseErrors(x, y int, img *image.RGBA, rErr, gErr, bErr, mul float64) color.Color {
	r, g, b, _ := img.At(x, y).RGBA()

	return color.RGBA{
		R: colours.ClampFloatToUint8(float64(r>>8) + rErr*mul),
		G: colours.ClampFloatToUint8(float64(g>>8) + gErr*mul),
		B: colours.ClampFloatToUint8(float64(b>>8) + bErr*mul),
		A: 255,
	}
}

func fsQuantisedErrors(c1, c2 color.Color) (float64, float64, float64) {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	rErr := float64(r2>>8) - float64(r1>>8)
	gErr := float64(g2>>8) - float64(g1>>8)
	bErr := float64(b2>>8) - float64(b1>>8)

	return rErr, gErr, bErr
}

// Bayer Dithering
func averageColourSpread(c color.Palette) float64 {
	var total = colours.Sqr(float64(len(c)))
	var dst = 0.0

	for i := range c {
		for j := range c {
			r1, g1, b1, _ := c[i].RGBA()
			r2, g2, b2, _ := c[j].RGBA()

			rDst := colours.Sqr(float64(r2>>8) - float64(r1>>8))
			gDst := colours.Sqr(float64(g2>>8) - float64(g1>>8))
			bDst := colours.Sqr(float64(b2>>8) - float64(b1>>8))

			dst += rDst + gDst + bDst
		}
	}

	return math.Sqrt(dst) / total
}

var bayerMatrix4x4 = [][]float64{
	{0, 8, 2, 10},
	{12, 4, 14, 6},
	{3, 11, 1, 9},
	{15, 7, 13, 5},
}

var bayerMatrix8x8 = [][]float64{
	{0, 32, 8, 40, 2, 34, 10, 42},
	{48, 16, 56, 24, 50, 18, 58, 26},
	{12, 44, 4, 36, 14, 46, 6, 38},
	{60, 28, 52, 20, 62, 30, 54, 22},
	{3, 35, 11, 43, 1, 33, 9, 41},
	{51, 19, 59, 27, 49, 17, 57, 25},
	{15, 47, 7, 39, 13, 45, 5, 37},
	{63, 31, 55, 23, 61, 29, 53, 21},
}

func bayerDither4x4(cimg *image.RGBA, c color.Palette) *image.RGBA {
	return bayerDitherWithOpts(cimg, c, bayerMatrix4x4)
}

func bayerDither8x8(cimg *image.RGBA, c color.Palette) *image.RGBA {
	return bayerDitherWithOpts(cimg, c, bayerMatrix8x8)
}

func bayerDitherWithOpts(cimg *image.RGBA, c color.Palette, matrix [][]float64) *image.RGBA {
	rowL := len(matrix[0])
	spread := averageColourSpread(c)
	bounds := cimg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			m := matrix[x%rowL][y%rowL]/16 - 0.5
			r, g, b, _ := cimg.At(x, y).RGBA()
			clr := c.Convert(color.RGBA{
				R: colours.ClampFloatToUint8(float64(r>>8) + spread*m),
				G: colours.ClampFloatToUint8(float64(g>>8) + spread*m),
				B: colours.ClampFloatToUint8(float64(b>>8) + spread*m),
				A: 255,
			})
			cimg.Set(x, y, clr)
		}
	}

	return cimg
}

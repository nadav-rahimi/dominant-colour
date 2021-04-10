package quantisers

import (
	"github.com/fiwippi/go-quantise/pkg/colours"
	"image"
	"image/color"
)

// Floyd-steinberg dithering https://en.wikipedia.org/wiki/Floyd%E2%80%93Steinberg_dithering

//
func diffuseErrors(x, y int, img *image.RGBA, rErr, gErr, bErr, mul float64) color.Color {
	r, g, b, _ := img.At(x, y).RGBA()

	return color.RGBA{
		R: colours.ClampUint8(float64(r>>8) + rErr*mul),
		G: colours.ClampUint8(float64(g>>8) + gErr*mul),
		B: colours.ClampUint8(float64(b>>8) + bErr*mul),
		A: 255,
	}
}

//
func quantisedErrors(c1, c2 color.Color) (float64, float64, float64) {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	rErr := float64(r2>>8) - float64(r1>>8)
	gErr := float64(g2>>8) - float64(g1>>8)
	bErr := float64(b2>>8) - float64(b1>>8)

	return rErr, gErr, bErr
}

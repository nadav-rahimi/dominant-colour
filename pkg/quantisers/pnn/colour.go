package pnn

import (
	"github.com/nadav-rahimi/dominant-colour/internal/quantisers/pnn"
	"image"
	"image/color"
)

// Returns a palette of "m" colours to best recreate the image from
func QuantiseColour(img image.Image, m int) color.Palette {
	return pnn.RGB.QuantiseColour(img, m)
}

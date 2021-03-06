package pnnlab

import (
	"github.com/fiwippi/go-quantise/internal/quantisers/pnn"
	"image"
	"image/color"
)

// Returns a palette of "m" colours to best recreate the image from
func QuantiseColour(img image.Image, m int) color.Palette {
	return pnn.LAB.QuantiseColour(img, m)
}

package pnnlab

import (
	"github.com/fiwippi/go-quantise/internal/quantisers/pnn"
	"image"
	"image/color"
)

// Returns "m" greyscale colours to best recreate the colour palette of the original image
func QuantiseGreyscale(img image.Image, m int) color.Palette {
	return pnn.QuantiseGreyscale(img, m)
}

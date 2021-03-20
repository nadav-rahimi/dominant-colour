package quantisers

import (
	"image"
)

// Linear Histogram which can represent a single colour
//channel or greyscale channel in the range 0-255
type LinearHistogram map[uint8]int

// Creates a linear histogram for the greyscale colour channel
func CreateGreyscaleHistogram(img image.Image) LinearHistogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(LinearHistogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Calculate the greyscale value (luminosity) of the pixels which are clamped to the range 0-255
			lum := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			pixels[lum]++
		}
	}

	return pixels
}

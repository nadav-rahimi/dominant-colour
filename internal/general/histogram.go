package general

import (
	"image"
)

// Linear Histogram used for single colour channels
type Histogram map[uint8]int

// Creates a linear histogram for the greyscale colour channel
func CreateGreyscaleHistogram(img image.Image) Histogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(Histogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r = r >> 8
			g = g >> 8
			b = b >> 8

			y := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			pixels[y]++
		}
	}

	return pixels
}

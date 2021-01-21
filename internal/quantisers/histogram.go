package quantisers

import (
	"image"
)

type Histogram map[uint8]int

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

func CreateRGBAHistogram(img image.Image) (rhist, ghist, bhist, ahist Histogram) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	r_pixels := make(Histogram)
	g_pixels := make(Histogram)
	b_pixels := make(Histogram)
	a_pixels := make(Histogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			r_pixels[uint8(r)]++
			g_pixels[uint8(g)]++
			b_pixels[uint8(b)]++
			a_pixels[uint8(a)]++

		}
	}

	return r_pixels, g_pixels, b_pixels, a_pixels
}

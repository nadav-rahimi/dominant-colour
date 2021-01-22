package pnn

import "image"

// Some code adapted from https://github.com/mcychan/PnnQuant.js/blob/master/src/pnnquant.js

type ColourHistogram map[uint32]*Node

func ARGBIndex(a, r, g, b uint32) uint32 {
	return (a&0xF0)<<8 | (r&0xF0)<<4 | (g & 0xF0) | (b >> 4)
}

func CreatePNNHistogram(img image.Image) ColourHistogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(ColourHistogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8
			index := ARGBIndex(a, r, g, b)

			// Create a node if it doesnt exist
			if pixels[index] == nil {
				pixels[index] = &Node{}
			}

			// Add the pixel to the bin
			pixels[index].A += float64(a)
			pixels[index].R += float64(r)
			pixels[index].G += float64(g)
			pixels[index].B += float64(b)
			pixels[index].N++
		}
	}

	return pixels
}

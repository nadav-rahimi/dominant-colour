package PNNLAB

import (
	"image"
)

// Histogram for PNN Nodes which are colour
type ColourHistogram map[uint32]*Node

// Converts a uint32 (A, R, G, B) colour into a 16 bit unique number
// The input parameters should be 8 bit, i.e. max of 255
func RGBIndex(r, g, b uint32) uint32 {
	return (r&0xF8)<<8 | (g&0xFC)<<3 | (b >> 3)
}

// Converts a uint32 (A, R, G, B) colour into a 16 bit unique number
// The input parameters should be 8 bit, i.e. max of 255
func ARGBIndex(a, r, g, b uint32) uint32 {
	return (a&0xF0)<<8 | (r&0xF0)<<4 | (g & 0xF0) | (b >> 4)
}

// Creates a coloured PNN Histogram
func CreatePNNLABHistogram(img image.Image) ColourHistogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(ColourHistogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r, g, b = r>>8, g>>8, b>>8

			// Converts the colour to a unique uint16 number
			index := ARGBIndex(r, g, b, 255)

			// Create a node if it doesnt exist
			if pixels[index] == nil {
				pixels[index] = &Node{}
			}

			// Add the pixel to the bin
			pixels[index].R += float64(r)
			pixels[index].G += float64(g)
			pixels[index].B += float64(b)
			pixels[index].N++
		}
	}

	return pixels
}

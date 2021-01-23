package PNN

import "image"

// Histogram for PNN Nodes which are colour
type ColourHistogram map[uint32]*Node

// Converts a uint32 (A, R, G, B) colour into a 16 bit unique number
// The input parameters should be 8 bit, i.e. max of 255
func ARGBIndex(a, r, g, b uint32) uint32 {
	return (a&0xF0)<<8 | (r&0xF0)<<4 | (g & 0xF0) | (b >> 4)
}

// Creates a coloured PNN Histogram
func CreatePNNHistogram(img image.Image) ColourHistogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(ColourHistogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			index := ARGBIndex(a>>8, r>>8, g>>8, b>>8)

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

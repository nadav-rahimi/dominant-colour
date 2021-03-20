package pnn

import "image"

// Histogram for PNN Nodes which are coloured
type Histogram map[uint32]*Node

// Takes uint32 RGBA colours and gives them a unique uint16 index value,
// this simplifies the colour space and speeds up computation without
// a noticeable loss in quality
func argbIndex(a, r, g, b uint32) uint32 {
	return (a&0xF0)<<8 | (r&0xF0)<<4 | (g & 0xF0) | (b >> 4)
}

// Creates a PNN Histogram
func CreatePNNHistogram(img image.Image) Histogram {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(Histogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			a, r, g, b = a>>8, r>>8, g>>8, b>>8

			// Get a uint16 number to use as an index for the colour
			index := argbIndex(a, r, g, b)

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

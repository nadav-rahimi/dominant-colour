package images

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// TODO explain if recreating from one colour then it splits at that colour otherwise image would be one whole colour
// TODO Support alpha channel
// TODO Why am I copying the original image onto the new image ?
// TODO move histograms from images folder to quantisers folder

//
func ImageFromColour(outputPath string, img image.Image, c []color.RGBA) error {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewRGBA(bounds)
	draw.Draw(cimg, bounds, image.Transparent, image.Point{}, draw.Over)

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			pixelColour := color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}

			if len(c) > 1 {
				maxDistance := math.MaxUint32
				bestColour := color.RGBA{}
				for _, v := range c {
					distance := colourDifference(pixelColour, v)
					if distance < maxDistance {
						maxDistance = distance
						bestColour = v
					}
				}

				cimg.Set(x, y, bestColour)
			} else if len(c) == 1 {
				cr := uint8(0)
				cg := uint8(0)
				cb := uint8(0)
				ca := uint8(0)
				if r >= uint32(c[0].R) {
					cr = 0xff
				}
				if g >= uint32(c[0].G) {
					cg = 0xff
				}
				if b >= uint32(c[0].B) {
					cb = 0xff
				}
				if a >= uint32(c[0].A) {
					ca = 0xff
				}

				cimg.Set(x, y, color.RGBA{cr, cg, cb, ca})
			} else {
				return errors.New("Colours must be specified")
			}

		}
	}

	err := SaveImage(outputPath, cimg)
	return err
}

// If the image is recreated at one colour, then it is
// recreated in Black and White and split at that colour
func ImageFromGreyscale(outputPath string, img image.Image, c []color.Gray) error {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewGray(bounds)
	draw.Draw(cimg, bounds, image.Transparent, image.Point{}, draw.Over)

	// Looping over Y first and X second is more likely to result
	// in better memory access patterns than X first and Y second.
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r = r >> 8
			g = g >> 8
			b = b >> 8

			greyscaleLevel := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

			if len(c) > 1 {
				maxDistance := 256
				bestColour := color.Gray{}
				for _, v := range c {
					distance := int(v.Y) - greyscaleLevel
					if distance < 0 {
						distance *= -1
					}
					if distance < maxDistance {
						maxDistance = distance
						bestColour.Y = v.Y
					}
				}

				cimg.Set(x, y, bestColour)
			} else if len(c) == 1 {
				if greyscaleLevel <= int(c[0].Y) {
					cimg.Set(x, y, color.Gray{0})
				} else {
					cimg.Set(x, y, color.Gray{255})
				}
			} else {
				return errors.New("Colours must be specified")
			}

		}
	}

	err := SaveImage(outputPath, cimg)
	return err
}

//
func colourDifference(c1, c2 color.RGBA) int {
	r_diff := math.Pow(float64(c2.R-c1.R), 2)
	g_diff := math.Pow(float64(c2.G-c1.G), 2)
	b_diff := math.Pow(float64(c2.B-c1.B), 2)
	a_diff := math.Pow(float64(c2.A-c1.A), 2)
	return int(math.Sqrt(r_diff + g_diff + b_diff + a_diff))
}

package quantisers

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"reflect"
)

// Recreates image from colour palette
func ImageFromPalette(img image.Image, c color.Palette) (image.Image, error) {
	if c == nil || len(c) < 1 {
		return nil, errors.New("Colours must be specified")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			pixelColour := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

			if len(c) > 1 {
				cimg.Set(x, y, c.Convert(pixelColour))
			} else if reflect.TypeOf(c[0]) == reflect.TypeOf(color.Gray{}) {
				Y := c[0].(color.Gray).Y
				greyscaleLevel := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

				if greyscaleLevel <= Y {
					cimg.Set(x, y, color.Gray{0})
				} else {
					cimg.Set(x, y, color.Gray{255})
				}
			} else if reflect.TypeOf(c[0]) == reflect.TypeOf(color.RGBA{}) {
				clr := c[0].(color.RGBA)
				cr := uint8(0)
				cg := uint8(0)
				cb := uint8(0)
				ca := uint8(0)
				if r >= uint32(clr.R) {
					cr = 0xff
				}
				if g >= uint32(clr.G) {
					cg = 0xff
				}
				if b >= uint32(clr.B) {
					cb = 0xff
				}
				if a >= uint32(clr.A) {
					ca = 0xff
				}
				cimg.Set(x, y, color.RGBA{cr, cg, cb, ca})
			} else {
				return nil, errors.New("Invalid colour palette")
			}
		}
	}

	return cimg, nil
}

// Rectangular image of all the colours in a palette stacked horizontally
// "ss" denotes the width/height of each colour square in pixels
func ColourPalette(c color.Palette, ss int) image.Image {
	numColours := len(c)
	img := image.NewRGBA(image.Rect(0, 0, ss*numColours, ss))

	for i, v := range c {
		uniform_colour := image.NewUniform(v)
		draw.Draw(img, image.Rect(ss*i, 0, ss*i+ss, ss), uniform_colour, image.Point{}, draw.Src)
	}

	return img
}

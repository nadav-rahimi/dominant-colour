package quantisers

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"reflect"
)

var (
	BLACK = color.Gray{Y: 0}
	WHITE = color.Gray{Y: 255}
)

// Recreates image from colour palette. If one greyscale colour is
// specified then the image is recreated in black and white with the
// split between them at the specified input colour
func ImageFromPalette(img image.Image, c color.Palette, dither bool) (image.Image, error) {
	if c == nil || len(c) < 1 {
		return nil, errors.New("Colours must be specified")
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewRGBA(bounds)
	draw.Draw(cimg, bounds, img, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			if len(c) > 1 {
				if dither {
					oldColour := cimg.At(x, y)
					newColour := c.Convert(oldColour)
					cimg.Set(x, y, newColour)

					rErr, gErr, bErr := quantisedErrors(newColour, oldColour)
					cimg.Set(x+1, y, diffuseErrors(x+1, y, cimg, rErr, gErr, bErr, 7.0/16))
					cimg.Set(x-1, y+1, diffuseErrors(x-1, y+1, cimg, rErr, gErr, bErr, 3.0/16))
					cimg.Set(x, y+1, diffuseErrors(x, y+1, cimg, rErr, gErr, bErr, 5.0/16))
					cimg.Set(x+1, y+1, diffuseErrors(x+1, y+1, cimg, rErr, gErr, bErr, 1.0/16))
				} else {
					cimg.Set(x, y, c.Convert(img.At(x, y)))
				}
			} else if reflect.TypeOf(c[0]) == reflect.TypeOf(color.Gray{}) {
				// Get the 8bit RGBA colours and calculate the greyscale equivalent
				r, g, b, _ := img.At(x, y).RGBA()
				r, g, b = r>>8, g>>8, b>>8
				greyscaleLevel := uint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
				// Get the Y threshold
				grayImage, ok := c[0].(color.Gray)
				if !ok {
					return nil, errors.New("Image could not be converted to greyscale")
				}

				if greyscaleLevel <= grayImage.Y {
					cimg.Set(x, y, BLACK)
				} else {
					cimg.Set(x, y, WHITE)
				}
			} else {
				return nil, errors.New("Invalid colour palette")
			}
		}
	}

	return cimg, nil
}

// Returns image of the colour palette, which each colour represented
// as a square. The size of the square in pixels is also specified.
func ColourPaletteImage(c color.Palette, size int) image.Image {
	numColours := len(c)
	img := image.NewRGBA(image.Rect(0, 0, size*numColours, size))

	for i, v := range c {
		uniform_colour := image.NewUniform(v)
		draw.Draw(img, image.Rect(size*i, 0, size*i+size, size), uniform_colour, image.Point{}, draw.Src)
	}

	return img
}

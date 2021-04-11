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
func ImageFromPalette(img image.Image, c color.Palette, ditherType DitherType) (image.Image, error) {
	if c == nil || len(c) < 1 {
		return nil, errors.New("Colours must be specified")
	}

	cimg := image.NewRGBA(img.Bounds())
	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Src)

	// Process one colour greyscale palettes
	if len(c) == 1 && reflect.TypeOf(c[0]) == reflect.TypeOf(color.Gray{}) {
		return noDitherSingle(cimg, c), nil
	}

	// Process multi colour palettes
	switch ditherType {
	case NoDither:
		return noDitherMulti(cimg, c), nil
	case FloydSteinberg:
		return floydSteinbergDither(cimg, c), nil
	case Bayer4x4:
		return bayerDither4x4(cimg, c), nil
	case Bayer8x8:
		return bayerDither8x8(cimg, c), nil
	default:
		return nil, errors.New("Invalid dither type")
	}
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

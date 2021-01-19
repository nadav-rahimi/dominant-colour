package images

import (
	"image"
	"image/color"
	"image/draw"
)

// TODO one function for this for greyscale and rgba called ColourPallette()

func DrawRectangleColour(path string, c []color.RGBA) error {
	numColours := len(c)
	img := image.NewRGBA(image.Rect(0, 0, 200*numColours, 200))

	for i, v := range c {
		uniform_colour := image.NewUniform(v)
		draw.Draw(img, image.Rect(200*i, 0, 200*i+200, 200), uniform_colour, image.ZP, draw.Src)
	}

	err := SaveImage(path, img)
	return err
}

func DrawRectangleGreyscale(path string, c []color.Gray) error {
	numColours := len(c)
	img := image.NewGray(image.Rect(0, 0, 200*numColours, 200))

	for i, v := range c {
		uniform_colour := image.NewUniform(v)
		draw.Draw(img, image.Rect(200*i, 0, 200*i+200, 200), uniform_colour, image.ZP, draw.Src)
	}

	err := SaveImage(path, img)
	return err
}

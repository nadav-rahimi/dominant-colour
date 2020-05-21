package dominant_colour

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
	"strings"
)

//
func readImage(path string) image.Image {
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	chkFatal(err, "Image could not be read")

	return m
}

func saveImage(path string, img image.Image) error {
	var encodeMethod int = 0
	if strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".jpg") {
		encodeMethod = 1
	} else if strings.HasSuffix(path, ".png") {
		encodeMethod = 2
	} else {
		return errors.New("File must be .jpeg/.jpg or .png")
	}

	// TODO validate the path
	toimg, err := os.Create(path)
	if err != nil {
		return err
	}
	defer toimg.Close()

	switch encodeMethod {
	case 1:
		jpeg.Encode(toimg, img, nil)
	case 2:
		png.Encode(toimg, img)
	}

	return nil
}

//
type histogram map[uint8]int

func createGreyscaleHistogram(path string) histogram {
	m := readImage(path)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	pixels := make(histogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

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

func createRGBAHistogram(path string) (rhist, ghist, bhist, ahist histogram) {
	m := readImage(path)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	r_pixels := make(histogram)
	g_pixels := make(histogram)
	b_pixels := make(histogram)
	a_pixels := make(histogram)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := m.At(x, y).RGBA()

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

//
func DrawSquareColour(c color.RGBA, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	uniform_colour := image.NewUniform(c)
	draw.Draw(img, image.Rect(0, 0, 200, 200), uniform_colour, image.ZP, draw.Src)

	err := saveImage(path, img)
	return err
}

func DrawSquareGrayscale(c color.Gray, path string) error {
	img := image.NewGray(image.Rect(0, 0, 200, 200))

	uniform_colour := image.NewUniform(c)
	draw.Draw(img, image.Rect(0, 0, 200, 200), uniform_colour, image.ZP, draw.Src)

	err := saveImage(path, img)
	return err
}

//
func RecreateImageFromColour(inputPath, outputPath string, c color.RGBA) error {
	m := readImage(inputPath)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewRGBA(m.Bounds())
	draw.Draw(cimg, m.Bounds(), m, image.Point{}, draw.Over)

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255
			r = r >> 8
			g = g >> 8
			b = b >> 8
			a = a >> 8

			cr := uint8(0)
			cg := uint8(0)
			cb := uint8(0)
			ca := uint8(0)
			if r >= uint32(c.R) {
				cr = 0xff
			}
			if g >= uint32(c.G) {
				cg = 0xff
			}
			if b >= uint32(c.B) {
				cb = 0xff
			}
			if a >= uint32(c.A) {
				ca = 0xff
			}

			cimg.Set(x, y, color.RGBA{cr, cg, cb, ca})

		}
	}

	err := saveImage(outputPath, cimg)
	return err
}

func RecreateImageFromGreyscale(inputPath, outputPath string, c color.Gray) error {
	m := readImage(inputPath)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewGray(m.Bounds())
	draw.Draw(cimg, m.Bounds(), m, image.Point{}, draw.Over)

	// Looping over Y first and X second is more likely to result
	// in better memory access patterns than X first and Y second.
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255 so 8 bits for grayscale
			r = r >> 8
			g = g >> 8
			b = b >> 8

			greyscaleLevel := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))

			if greyscaleLevel <= int(c.Y) {
				cimg.Set(x, y, color.Gray{0})
			} else {
				cimg.Set(x, y, color.Gray{255})
			}
		}
	}

	err := saveImage(outputPath, cimg)
	return err
}

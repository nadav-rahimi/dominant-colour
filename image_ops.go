package dominant_colour

import (
	"fmt"
	golang_sets "github.com/nadav-rahimi/golang-sets"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

// Struct to hold the RGB colour data the program
// calculates, values should be between 0 and 255
// inclusive
type RGB struct {
	R float64
	G float64
	B float64
}

// Draws a rectangle of 200x200 squares of the colours
// input into the function
func DrawRectangle(colours []*RGB, path string) {
	numColours := len(colours)
	img := image.NewRGBA(image.Rect(0, 0, 200*numColours, 200))

	for i := 0; i < len(colours); i++ {
		r := uint8(colours[i].R)
		g := uint8(colours[i].G)
		b := uint8(colours[i].B)
		c := image.NewUniform(color.RGBA{r, g, b, 0xff})
		draw.Draw(img, image.Rect(200*i, 0, 200*i+200, 200), c, image.ZP, draw.Src)
	}

	// TODO validate the path
	toimg, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}

// Calculates the difference between two RGB colours
// by treating them as vectors in 3D space
func distanceBetween(c1, c2 *RGB) float64 {
	r_diff := math.Pow(c2.R-c1.R, 2)
	g_diff := math.Pow(c2.G-c1.G, 2)
	b_diff := math.Pow(c2.B-c1.B, 2)
	return math.Sqrt(r_diff + g_diff + b_diff)
}

// Reads JPEG and PNG images into a set of RGB vectors
func Img2pixelset(path string) *golang_sets.Set {
	//Decode the JPEG data.
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Looping over Y first and X second is more likely to result
	// in better memory access patterns than X first and Y second.
	pixels := golang_sets.NewSet()
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255
			r = r >> 8
			g = g >> 8
			b = b >> 8

			pixels.Add(mat.NewVecDense(3, []float64{float64(r), float64(g), float64(b)}))
		}
	}

	return pixels
}

// Recreates a given image as best as possible using the given RGB colours
func RecreateImage(imgpath string, colours []*RGB) {
	reader, err := os.Open(imgpath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	bounds := m.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cimg := image.NewRGBA(m.Bounds())
	draw.Draw(cimg, m.Bounds(), m, image.Point{}, draw.Over)

	// Looping over Y first and X second is more likely to result
	// in better memory access patterns than X first and Y second.
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			r, g, b, _ := m.At(x, y).RGBA()

			// Convert rgb values to be in range 0-255
			r = r >> 8
			g = g >> 8
			b = b >> 8

			c := &RGB{
				R: float64(r),
				G: float64(g),
				B: float64(b),
			}

			difference := float64(442)
			c_sub := &RGB{
				R: 0,
				G: 0,
				B: 0,
			}
			for i := range colours {
				temp := distanceBetween(c, colours[i])
				if temp < difference {
					difference = temp
					c_sub = colours[i]
				}
			}

			cimg.Set(x, y, color.RGBA{uint8(c_sub.R), uint8(c_sub.G), uint8(c_sub.B), 255})
		}
	}

	output := fmt.Sprintf("%s_render.jpeg", imgpath[:len(imgpath)-4])
	outFile, _ := os.Create(output)
	defer outFile.Close()
	jpeg.Encode(outFile, cimg, nil)
}

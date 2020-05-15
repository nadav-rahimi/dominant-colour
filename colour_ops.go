package dominant_colour

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
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
func DrawRectangle(colours []*RGB) {
	numColours := len(colours)
	img := image.NewRGBA(image.Rect(0, 0, 200*numColours, 200))

	for i := 0; i < len(colours); i++ {
		r := uint8(colours[i].R)
		g := uint8(colours[i].G)
		b := uint8(colours[i].B)
		c := image.NewUniform(color.RGBA{r, g, b, 0xff})
		draw.Draw(img, image.Rect(200*i, 0, 200*i + 200, 200), c,image.ZP, draw.Src)
	}

	toimg, err := os.Create("dominantcolours.png")
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

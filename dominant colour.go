package dominant_colour

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/nadav-rahimi/golang-sets"
	"gonum.org/v1/gonum/mat"
)

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

// Find "n" most dominant colours with the path to the image given
// Uses a custom binary tree
func FindDominantColoursBT(path string, n int) []*RGB {
	root := newNode(Img2pixelset(path))
	root.calculate_mean_and_covariance()

	for i := 0; i < n; i++ {
		fmt.Printf("WORKING ON ITERATION: %v\n", i)
		node := root.find_max_eigenvector()
		node.calculate_mean_and_covariance()
		node.partition_node()
	}

	colours := make([]*RGB, 0, n)
	leaves := root.get_leaves()
	for i := range leaves {
		r := leaves[i].qn.At(0, 0)
		g := leaves[i].qn.At(1, 0)
		b := leaves[i].qn.At(2, 0)
		colours = append(colours, &RGB{r, g, b})
	}

	return colours
}

// If you want to reuse the same pixel set multiple times you can read
// the data into a set using the Img3pixelset function and then pass it
// into this function which makes a copy of the set before calculating
// the dominant colours leaving the original set untouched
func FindDominantColoursBTFromSet(s *golang_sets.Set, n int) []*RGB {
	root := newNode(s)
	root.calculate_mean_and_covariance()

	for i := 0; i < n; i++ {
		fmt.Printf("WORKING ON ITERATION: %v\n", i)
		node := root.find_max_eigenvector()
		node.calculate_mean_and_covariance()
		node.partition_node()
	}

	colours := make([]*RGB, 0, n)
	leaves := root.get_leaves()
	for i := range leaves {
		r := leaves[i].qn.At(0, 0)
		g := leaves[i].qn.At(1, 0)
		b := leaves[i].qn.At(2, 0)
		colours = append(colours, &RGB{r, g, b})
	}

	return colours
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

// Recreates a given image as best as possible using the given RGB colours
// With name appended to the end, name is int
func RecreateImageAppendName(imgpath string, colours []*RGB, n int) {
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

	output := fmt.Sprintf("%s_render_%v.jpeg", imgpath[:len(imgpath)-4], n)
	outFile, _ := os.Create(output)
	defer outFile.Close()
	jpeg.Encode(outFile, cimg, nil)
}

// Please dont use this method yet
func EfficientMultipleImageGen(path string, n int) {
	root := newNode(Img2pixelset(path))
	root.calculate_mean_and_covariance()

	for i := 0; i < n; i++ {
		fmt.Printf("%s, WORKING ON ITERATION: %v\n", path, i)
		node := root.find_max_eigenvector()
		node.calculate_mean_and_covariance()
		node.partition_node()

		colours := make([]*RGB, 0, n)
		leaves := root.get_leaves()
		for i := range leaves {
			r := leaves[i].qn.At(0, 0)
			g := leaves[i].qn.At(1, 0)
			b := leaves[i].qn.At(2, 0)
			colours = append(colours, &RGB{r, g, b})
		}

		RecreateImageAppendName(path, colours, i)
	}
}

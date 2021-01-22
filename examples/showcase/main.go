package main

import (
	"flag"
	"fmt"
	"github.com/nadav-rahimi/dominant-colour/examples/showcase/pkg/images"
	"github.com/nadav-rahimi/dominant-colour/pkg/quantisers"
	LMQ "github.com/nadav-rahimi/dominant-colour/pkg/quantisers/lmq"
	Otsu "github.com/nadav-rahimi/dominant-colour/pkg/quantisers/otsu"
	PNN "github.com/nadav-rahimi/dominant-colour/pkg/quantisers/pnn"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
)

// TODO add readme and .gitkeep and a go mod file

var (
	multi        *int
	err          error
	useRandom    *bool
	quantisedImg image.Image
	palette      image.Image
	img          image.Image
	colours      color.Palette
)

func main() {
	useRandom = flag.Bool("random", false, "Whether to use a random image instead of the supplied one")
	multi = flag.Int("multi", 6, "Number of colours for multi tone quantisations")
	flag.Parse()

	img, err = images.ReadImage("graffiti.jpg")
	if *useRandom {
		img = randomImage()
	}

	// Greyscale Image Single Tone
	colours, err = PNN.Colour(img, *multi)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("pnn-colour-multi.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("pnn-colour-multi-palette.png", palette, images.BestSpeed))

	//OtsuExample()
	//LMQExample()
	//PNNExample()
}

func OtsuExample() {
	fmt.Println("Creating Otsu...")

	// Greyscale Image
	colours, err = Otsu.Greyscale(img)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("otsu-greyscale-single.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("otsu-greyscale-single-palette.png", palette, images.BestSpeed))

	// Colour Image
	colours, err = Otsu.Colour(img)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("otsu-colour-single.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("otsu-colour-single-palette.png", palette, images.BestSpeed))

	fmt.Println("Finished Otsu")
}

func PNNExample() {
	fmt.Println("Creating PNN...")

	// Greyscale Image Single Tone
	colours, err = PNN.Greyscale(img, 1)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("pnn-greyscale-single.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("pnn-greyscale-single-palette.png", palette, images.BestSpeed))

	// Greyscale Image Multi Tone
	colours, err = PNN.Greyscale(img, *multi)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("pnn-greyscale-multi.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("pnn-greyscale-multi-palette.png", palette, images.BestSpeed))

	fmt.Println("Finished PNN")
}

func LMQExample() {
	fmt.Println("Creating LMQ...")

	// Greyscale Image Single Tone
	colours, err = LMQ.Greyscale(img, 1)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("lmq-greyscale-single.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("lmq-greyscale-single-palette.png", palette, images.BestSpeed))

	// Greyscale Image Multi Tone
	colours, err = LMQ.Greyscale(img, *multi)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("lmq-greyscale-multi.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("lmq-greyscale-multi-palette.png", palette, images.BestSpeed))

	// Colour Image
	colours, err = LMQ.Colour(img)
	logErr(err)
	quantisedImg, err = quantisers.ImageFromPalette(img, colours)
	logErr(err)
	palette = quantisers.ColourPalette(colours, 200)
	logErr(images.SaveImage("lmq-colour-single.png", quantisedImg, images.BestSpeed))
	logErr(images.SaveImage("lmq-colour-single-palette.png", palette, images.BestSpeed))

	fmt.Println("Finished LMQ")
}

func randomImage() image.Image {
	fmt.Println("Using Random...")

	resp, err := http.Get("https://picsum.photos/3000/2000")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	img, err = jpeg.Decode(resp.Body)

	err = images.SaveImage("random.png", img, images.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}

	// Save greyscale version of image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	cimg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			oldPixel := img.At(x, y)
			pixel := color.GrayModel.Convert(oldPixel)
			cimg.Set(x, y, pixel)
		}
	}
	err = images.SaveImage("random-grey.png", cimg, images.BestSpeed)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

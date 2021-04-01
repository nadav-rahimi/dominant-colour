package main

import (
	"flag"
	"fmt"
	"github.com/fiwippi/go-quantise/pkg/quantisers"
	"github.com/fiwippi/go-quantise/pkg/quantisers/lmq"
	"github.com/fiwippi/go-quantise/pkg/quantisers/otsu"
	"github.com/fiwippi/go-quantise/pkg/quantisers/pnn"
	"github.com/fiwippi/go-quantise/pkg/quantisers/pnnlab"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"showcase/pkg/images"
)

var err error
var img image.Image
var paletteSize int

func main() {
	var randomImg bool

	flag.BoolVar(&randomImg, "random", false, "Whether to use a random image instead of the supplied one")
	flag.IntVar(&paletteSize, "palette-size", 7, "The size of the palette you want to generate")
	flag.Parse()

	if randomImg {
		img, err = randomImage()
	} else {
		img, err = images.ReadImage("fish.jpg")
	}
	if err != nil {
		log.Fatal(err)
	}

	OtsuExample()
	LMQExample()
	PNNExample()
	PNNLABExample()
}

func OtsuExample() {
	fmt.Println("Creating Otsu...")

	colours := otsu.QuantiseGreyscale(img)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours)
	palette := quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("otsu-grey.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("otsu-grey-palette.jpg", palette, images.BestSpeed)

	fmt.Println("Finished Otsu")
}

func LMQExample() {
	fmt.Println("Creating LMQ...")

	// Greyscale Image Single Tone
	colours := lmq.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours)
	palette := quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("lmq-grey-single.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("lmq-grey-single-palette.jpg", palette, images.BestSpeed)

	// Greyscale Image Multi Tone
	colours = lmq.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours)
	palette = quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("lmq-grey-multi.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("lmq-grey-multi-palette.jpg", palette, images.BestSpeed)

	fmt.Println("Finished LMQ")
}

func PNNExample() {
	fmt.Println("Creating PNN...")

	// Greyscale Image Single Tone
	colours := pnn.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours)
	palette := quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("pnn-grey-single.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("pnn-grey-single-palette.jpg", palette, images.BestSpeed)

	// Greyscale Image Multi Tone
	colours = pnn.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours)
	palette = quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("pnn-grey-multi.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("pnn-grey-multi-palette.jpg", palette, images.BestSpeed)

	// Colour Image Multi Tone
	colours = pnn.QuantiseColour(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours)
	palette = quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("pnn-colour-multi.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("pnn-colour-multi-palette.jpg", palette, images.BestSpeed)

	fmt.Println("Finished PNN")
}

func PNNLABExample() {
	fmt.Println("Creating PNN LAB...")

	// Doesn't do greyscale because same as PNN

	// Colour Image Multi Tone
	colours := pnnlab.QuantiseColour(img, paletteSize)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours)
	palette := quantisers.ColourPaletteImage(colours, 200)
	images.SaveImage("pnnlab-colour-multi.jpg", quantisedImg, images.BestSpeed)
	images.SaveImage("pnnlab-colour-multi-palette.jpg", palette, images.BestSpeed)

	fmt.Println("Finished PNN LAB")
}

func randomImage() (image.Image, error) {
	resp, err := http.Get("https://picsum.photos/3000/2000")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err = jpeg.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// Save coloured version of the image
	err = images.SaveImage("random.png", img, images.BestSpeed)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return img, nil
}

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
	"log"
)

var err error
var randomImg bool
var img image.Image
var paletteSize int

func main() {
	flag.BoolVar(&randomImg, "random", false, "Whether to use a random image instead of the supplied one")
	flag.IntVar(&paletteSize, "palette-size", 7, "The size of the palette you want to generate")
	flag.Parse()

	if randomImg {
		img, err = randomImage()
	} else {
		img, err = ReadImage("fish.jpg")
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
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, false)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("otsu-grey.jpg", quantisedImg)
	SaveJPEG("otsu-grey-palette.jpg", palette)

	fmt.Println("Finished Otsu")
}

func LMQExample() {
	fmt.Println("Creating LMQ...")

	// Greyscale Image Single Tone
	colours := lmq.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, false)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("lmq-grey-single.jpg", quantisedImg)
	SaveJPEG("lmq-grey-single-palette.jpg", palette)

	// Greyscale Image Multi Tone
	colours = lmq.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, false)
	quantisedImgDithered, _ := quantisers.ImageFromPalette(img, colours, true)
	palette = quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("lmq-grey-multi.jpg", quantisedImg)
	SaveJPEG("lmq-grey-multi-dithered.jpg", quantisedImgDithered)
	SaveJPEG("lmq-grey-multi-palette.jpg", palette)

	fmt.Println("Finished LMQ")
}

func PNNExample() {
	fmt.Println("Creating PNN...")

	// Greyscale Image Single Tone
	colours := pnn.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, false)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnn-grey-single.jpg", quantisedImg)
	SaveJPEG("pnn-grey-single-palette.jpg", palette)

	// Greyscale Image Multi Tone
	colours = pnn.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, false)
	quantisedImgDithered, _ := quantisers.ImageFromPalette(img, colours, true)
	palette = quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnn-grey-multi.jpg", quantisedImg)
	SaveJPEG("pnn-grey-multi-dithered.jpg", quantisedImgDithered)
	SaveJPEG("pnn-grey-multi-palette.jpg", palette)

	// Colour Image Multi Tone
	colours = pnn.QuantiseColour(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, false)
	quantisedImgDithered, _ = quantisers.ImageFromPalette(img, colours, true)
	palette = quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnn-colour-multi.jpg", quantisedImg)
	SaveJPEG("pnn-colour-multi-dithered.jpg", quantisedImgDithered)
	SaveJPEG("pnn-colour-multi-palette.jpg", palette)

	fmt.Println("Finished PNN")
}

func PNNLABExample() {
	fmt.Println("Creating PNN LAB...")

	// Doesn't do greyscale because same as PNN

	// Colour Image Multi Tone
	colours := pnnlab.QuantiseColour(img, paletteSize)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, false)
	quantisedImgDithered, _ := quantisers.ImageFromPalette(img, colours, true)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnnlab-colour-multi.jpg", quantisedImg)
	SaveJPEG("pnnlab-colour-multi-dithered.jpg", quantisedImgDithered)
	SaveJPEG("pnnlab-colour-multi-palette.jpg", palette)

	fmt.Println("Finished PNN LAB")
}

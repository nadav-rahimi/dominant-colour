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
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	SaveJPEG("otsu-grey.jpg", quantisedImg)

	fmt.Println("Finished Otsu")
}

func LMQExample() {
	fmt.Println("Creating LMQ...")

	// Greyscale Image Single Tone
	colours := lmq.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	SaveJPEG("lmq-grey-single.jpg", quantisedImg)

	// Greyscale Image Multi Tone
	colours = lmq.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("lmq-grey-multi.jpg", quantisedImg)
	SaveJPEG("lmq-grey-multi-palette.jpg", palette)

	fmt.Println("Finished LMQ")
}

func PNNExample() {
	fmt.Println("Creating PNN...")

	// Greyscale Image Single Tone
	colours := pnn.QuantiseGreyscale(img, 1)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	SaveJPEG("pnn-grey-single.jpg", quantisedImg)

	// Greyscale Image Multi Tone
	colours = pnn.QuantiseGreyscale(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnn-grey-multi.jpg", quantisedImg)
	SaveJPEG("pnn-grey-multi-palette.jpg", palette)

	// Colour Image Multi Tone
	colours = pnn.QuantiseColour(img, paletteSize)
	quantisedImg, _ = quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	ditheredFS, _ := quantisers.ImageFromPalette(img, colours, quantisers.FloydSteinberg)
	ditheredFSS, _ := quantisers.ImageFromPalette(img, colours, quantisers.FloydSteinbergSerpentine)
	ditheredB2x2, _ := quantisers.ImageFromPalette(img, colours, quantisers.Bayer2x2)
	ditheredB4x4, _ := quantisers.ImageFromPalette(img, colours, quantisers.Bayer4x4)
	ditheredB8x8, _ := quantisers.ImageFromPalette(img, colours, quantisers.Bayer8x8)
	palette = quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnn-colour-multi.jpg", quantisedImg)
	SaveJPEG("pnn-colour-multi-dithered-floydsteinberg.jpg", ditheredFS)
	SaveJPEG("pnn-colour-multi-dithered-floydsteinbergserpentine.jpg", ditheredFSS)
	SaveJPEG("pnn-colour-multi-dithered-bayer2x2.jpg", ditheredB2x2)
	SaveJPEG("pnn-colour-multi-dithered-bayer4x4.jpg", ditheredB4x4)
	SaveJPEG("pnn-colour-multi-dithered-bayer8x8.jpg", ditheredB8x8)
	SaveJPEG("pnn-colour-multi-palette.jpg", palette)

	fmt.Println("Finished PNN")
}

func PNNLABExample() {
	fmt.Println("Creating PNN LAB...")

	// Doesn't do greyscale because same as PNN

	// Colour Image Multi Tone
	colours := pnnlab.QuantiseColour(img, paletteSize)
	quantisedImg, _ := quantisers.ImageFromPalette(img, colours, quantisers.NoDither)
	palette := quantisers.ColourPaletteImage(colours, 200)
	SaveJPEG("pnnlab-colour-multi.jpg", quantisedImg)
	SaveJPEG("pnnlab-colour-multi-palette.jpg", palette)

	fmt.Println("Finished PNN LAB")
}

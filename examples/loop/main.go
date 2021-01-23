package main

import (
	"flag"
	"github.com/cheggaaa/pb/v3"
	"github.com/kettek/apng"
	"github.com/nadav-rahimi/dominant-colour/pkg/quantisers"
	"image"
	"image/color"
	"log"
	"loop/pkg/images"
	"os"
)

var maxImages int
var imagePath string
var loopName string

// TODO option to reverse loop

func main() {
	flag.StringVar(&imagePath, "img", "graffiti.jpg", "Path to the image to loop")
	flag.StringVar(&loopName, "name", "loop.png", "Name of the output apng file")
	flag.IntVar(&maxImages, "colours", 30, "Max number of quantised colours to loop to")
	flag.Parse()

	// Read the image
	var img, err = images.ReadImage(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create the apng encoder
	a := apng.APNG{
		Frames:    make([]apng.Frame, maxImages-1),
		LoopCount: 0,
	}
	out, err := os.Create(loopName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Creating each frame and saving the frame name to a list
	var colours color.Palette
	var quantisedImg image.Image

	// Creating the quantiser
	PNN := quantisers.NewPNNQuantiser()

	bar := pb.StartNew(maxImages)
	for i := 2; i <= maxImages; i++ {
		colours, err = PNN.Colour(img, i)
		if err != nil {
			log.Fatal(err)
		}
		quantisedImg, err = quantisers.ImageFromPalette(img, colours)
		if err != nil {
			log.Fatal(err)
		}
		a.Frames[i-2].Image = quantisedImg
		bar.Increment()
	}
	// Write APNG to our output file
	apng.Encode(out, a)
	// Signal program is done
	bar.Increment()
	bar.Finish()
}

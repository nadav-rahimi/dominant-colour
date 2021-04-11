package main

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"os"
)

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
	err = SaveJPEG("random.jpg", img)
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
	err = SaveJPEG("random-grey.jpg", cimg)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ReadImage(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func SaveJPEG(path string, img image.Image) error {
	toimg, err := os.Create(path)
	if err != nil {
		return err
	}
	defer toimg.Close()

	if err = jpeg.Encode(toimg, img, nil); err != nil {
		return err
	}

	return nil
}

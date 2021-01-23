package Otsu

import (
	"errors"
	q "github.com/nadav-rahimi/dominant-colour/internal/general"
	"image"
	"image/color"
	"sync"
)

type Otsu struct{}

// Returns one greyscale colour which best represents the threshold for splitting
// the image into black and white. Otsu only supports m = 1 but still takes in the
// m parameter to comply with the Quantiser interface
func (otsu Otsu) Greyscale(img image.Image, m int) (color.Palette, error) {
	if m > 1 {
		return nil, errors.New("Otsu does not support greyscale quantisation with m > 1")
	}

	histogram := q.CreateGreyscaleHistogram(img)
	ychan := make(chan int)
	go quantise(histogram, ychan, nil)
	T := <-ychan

	return color.Palette{color.Gray{uint8(T)}}, nil
}

// Implemented in accordance with the quantiser interface but does not support colour
// quantisation, only greyscale quantisation
func (otsu Otsu) Colour(img image.Image, m int) (color.Palette, error) {
	return nil, errors.New("Otsu does not support colour quantisation")
}

// Quantises a given histogram
func quantise(hist q.Histogram, c chan int, wg *sync.WaitGroup) {
	xMax := 256
	P := make([]int, xMax)
	S := make([]int, xMax)
	P[0] = 0
	S[0] = 0
	H := make([][]int, xMax)
	for i := range H {
		H[i] = make([]int, xMax)
	}

	for v := 0; v < xMax-1; v++ {
		P[v+1] = P[v] + hist[uint8(v+1)]
		S[v+1] = S[v] + (v+1)*hist[uint8(v+1)]
	}

	for v := 1; v < xMax; v++ {
		for u := 1; u < v; u++ {
			Sp := getSum(P, u, v)
			Ss := getSum(S, u, v)
			if getSum(P, u, v) == 0 {
				H[u][v] = 0
				continue
			}

			H[u][v] = (Ss / Sp) * Ss
		}
	}

	maxVariation := 0
	T := 0
	for t := 0; t < xMax-1; t++ {
		variation := H[1][t] + H[t+1][xMax-1]
		if variation > maxVariation {
			maxVariation = variation
			T = t
		}
	}

	if wg != nil {
		wg.Done()
	}

	c <- T
}

func getSum(sum []int, u, v int) int {
	return sum[v] - sum[u-1]
}

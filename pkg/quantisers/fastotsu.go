package quantisers

import (
	"errors"
	"github.com/nadav-rahimi/dominant-colour/pkg/images"
	"image"
	"image/color"
	"sync"
)

// TODO fix so no need for curly brackets when calling quantisers

type Otsu struct{}

func (o Otsu) Greyscale(img image.Image, m int) ([]color.Gray, error) {
	if m > 1 {
		return nil, errors.New("For 'm > 1' this function is not implemented due to memory constraints")
	}

	histogram := images.CreateGreyscaleHistogram(img)
	ychan := make(chan int)
	go o.bilevelThreshold(histogram, ychan, nil)
	T := <-ychan

	return []color.Gray{color.Gray{uint8(T)}}, nil
}

func (o Otsu) Colour(img image.Image, m int) ([]color.RGBA, error) {
	if m > 1 {
		return nil, errors.New("For 'm > 1' this function is not implemented due to memory constraints")
	}

	rhist, ghist, bhist, ahist := images.CreateRGBAHistogram(img)

	rchan := make(chan int)
	gchan := make(chan int)
	bchan := make(chan int)
	achan := make(chan int)

	var wg sync.WaitGroup
	wg.Add(4)

	go o.bilevelThreshold(rhist, rchan, &wg)
	go o.bilevelThreshold(ghist, gchan, &wg)
	go o.bilevelThreshold(bhist, bchan, &wg)
	go o.bilevelThreshold(ahist, achan, &wg)

	Tr := <-rchan
	Tg := <-gchan
	Tb := <-bchan
	Ta := <-achan

	wg.Wait()

	return []color.RGBA{color.RGBA{uint8(Tr), uint8(Tg), uint8(Tb), uint8(Ta)}}, nil
}

func (o Otsu) bilevelThreshold(hist images.GreyscaleHistogram, c chan int, wg *sync.WaitGroup) {
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
			Sp := o.getSum(P, u, v)
			Ss := o.getSum(S, u, v)
			if o.getSum(P, u, v) == 0 {
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

func (o Otsu) getSum(sum []int, u, v int) int {
	return sum[v] - sum[u-1]
}

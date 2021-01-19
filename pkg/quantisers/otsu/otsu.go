package Otsu

import (
	q "github.com/nadav-rahimi/dominant-colour/internal/quantisers"
	"image"
	"image/color"
	"sync"
)

func Greyscale(img image.Image) (color.Palette, error) {
	histogram := q.CreateGreyscaleHistogram(img)
	ychan := make(chan int)
	go bilevelThreshold(histogram, ychan, nil)
	T := <-ychan

	return color.Palette{color.Gray{uint8(T)}}, nil
}

func Colour(img image.Image) (color.Palette, error) {
	rhist, ghist, bhist, ahist := q.CreateRGBAHistogram(img)

	rchan := make(chan int)
	gchan := make(chan int)
	bchan := make(chan int)
	achan := make(chan int)

	var wg sync.WaitGroup
	wg.Add(4)

	go bilevelThreshold(rhist, rchan, &wg)
	go bilevelThreshold(ghist, gchan, &wg)
	go bilevelThreshold(bhist, bchan, &wg)
	go bilevelThreshold(ahist, achan, &wg)

	Tr := <-rchan
	Tg := <-gchan
	Tb := <-bchan
	Ta := <-achan

	wg.Wait()

	return color.Palette{color.RGBA{uint8(Tr), uint8(Tg), uint8(Tb), uint8(Ta)}}, nil
}

func bilevelThreshold(hist q.Histogram, c chan int, wg *sync.WaitGroup) {
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

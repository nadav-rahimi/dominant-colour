package dominant_colour

import (
	"errors"
	"image/color"
	"sync"
)

type OtsuStruct struct{}

func (OtsuStruct) DominantGreyscaleValue(input string) (color.Gray, error) {
	histogram := createGreyscaleHistogram(input)
	ychan := make(chan int)
	go OtsuStruct{}.bilevelThreshold(histogram, ychan, nil)
	T := <-ychan

	return color.Gray{uint8(T)}, nil
}

func (OtsuStruct) DominantGreyscaleValues(input string, m int) ([]color.Gray, error) {
	return nil, errors.New("This function is not implemented due to memory constraints")
}

func (OtsuStruct) DominantColourValue(input string) (color.RGBA, error) {
	rhist, ghist, bhist, ahist := createRGBAHistogram(input)

	rchan := make(chan int)
	gchan := make(chan int)
	bchan := make(chan int)
	achan := make(chan int)

	var wg sync.WaitGroup
	wg.Add(4)

	go OtsuStruct{}.bilevelThreshold(rhist, rchan, &wg)
	go OtsuStruct{}.bilevelThreshold(ghist, gchan, &wg)
	go OtsuStruct{}.bilevelThreshold(bhist, bchan, &wg)
	go OtsuStruct{}.bilevelThreshold(ahist, achan, &wg)

	Tr := <-rchan
	Tg := <-gchan
	Tb := <-bchan
	Ta := <-achan

	wg.Wait()

	return color.RGBA{uint8(Tr), uint8(Tg), uint8(Tb), uint8(Ta)}, nil
}

func (OtsuStruct) DominantColourValues(input string, m int) ([]color.RGBA, error) {
	return nil, errors.New("This function is not implemented due to memory constraints")
}

func (OtsuStruct) bilevelThreshold(hist histogram, c chan int, wg *sync.WaitGroup) {
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
			Sp := OtsuStruct{}.getSum(P, u, v)
			Ss := OtsuStruct{}.getSum(S, u, v)
			if (OtsuStruct{}.getSum(P, u, v) == 0) {
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

func (OtsuStruct) getSum(sum []int, u, v int) int {
	return sum[v] - sum[u-1]
}

package quantisers

import (
	"errors"
	"github.com/nadav-rahimi/dominant-colour/pkg/images"
	"image"
	"image/color"
	"sync"
)

// TODO Correct error message

type LMQ struct{}

func (lmq LMQ) Greyscale(img image.Image, m int) ([]color.Gray, error) {
	histogram := images.CreateGreyscaleHistogram(img)
	ychan := make(chan []int)
	go lmq.multilevelThreshold(histogram, m, ychan, nil)
	T := <-ychan

	colours := make([]color.Gray, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours, nil
}

func (lmq LMQ) Colour(img image.Image, m int) ([]color.RGBA, error) {
	if m > 1 {
		return nil, errors.New("For 'm > 1' this function is not implemented due to memory constraints")
	}

	rhist, ghist, bhist, ahist := images.CreateRGBAHistogram(img)

	rchan := make(chan []int)
	gchan := make(chan []int)
	bchan := make(chan []int)
	achan := make(chan []int)

	var wg sync.WaitGroup
	wg.Add(4)

	go lmq.multilevelThreshold(rhist, m, rchan, &wg)
	go lmq.multilevelThreshold(ghist, m, gchan, &wg)
	go lmq.multilevelThreshold(bhist, m, bchan, &wg)
	go lmq.multilevelThreshold(ahist, m, achan, &wg)

	Tr := <-rchan
	Tg := <-gchan
	Tb := <-bchan
	Ta := <-achan

	wg.Wait()

	colours := make([]color.RGBA, len(Tr))
	for i := range colours {
		colours[i] = color.RGBA{uint8(Tr[i]), uint8(Tg[i]), uint8(Tb[i]), uint8(Ta[i])}
	}

	return colours, nil
}

func (lmq LMQ) multilevelThreshold(hist images.GreyscaleHistogram, m int, c chan []int, wg *sync.WaitGroup) {
	xMax := 255
	xMin := 0

	// Calculate the initial threshold values
	T := make([]uint8, m+1)
	for i := 0; i <= m; i++ {
		T[i] = uint8(xMin + (i*(xMax-xMin))/m)
	}
	// Initialising the segment map
	segments := make(map[int]images.GreyscaleHistogram)
	for i := 1; i <= m; i++ {
		segments[i] = images.GreyscaleHistogram{}
	}
	// Initialising the averages for each segment
	averages := make(map[int]int)
	// Initialising the slice for the old threshold history
	oldT := make([]uint8, len(T))

	// Calculating the thresholds
	for {
		copy(oldT, T)

		// Segments the pixels of the image into thresholds based on the linearHistogram
		for k, v := range hist {
			if k == 0 {
				segments[1][0] = v
			}

			// Checking in general
			for i := 1; i <= m; i++ {
				if T[i-1] < k && k <= T[i] {
					segments[i][k] = v
				}
			}
		}

		// Calculating the segment averages
		for i := 1; i <= m; i++ {
			averages[i] = mean(segments[i])
		}

		// Recalculating the thresholds
		for i := 1; i <= m-1; i++ {
			T[i] = uint8((averages[i] + averages[i+1]) / 2)
		}

		// If the old threshold is equal to the new threshold we are done
		if equal(oldT, T) {
			break
		}

	}

	//fmt.Println(T)

	averagesToReturn := make([]int, 0, m)
	for i := 1; i <= m; i++ {
		averagesToReturn = append(averagesToReturn, averages[i])
	}

	if wg != nil {
		wg.Done()
	}

	c <- averagesToReturn
}

func mean(h images.GreyscaleHistogram) int {
	sum := 0
	total := 0

	for k, v := range h {
		sum += int(k) * v
		total += v
	}

	if total == 0 {
		return 0
	}
	return sum / total
}

func equal(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

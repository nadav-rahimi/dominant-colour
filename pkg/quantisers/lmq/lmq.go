package LMQ

import (
	q "github.com/nadav-rahimi/dominant-colour/internal/quantisers"
	"image"
	"image/color"
	"sync"
)

func Greyscale(img image.Image, m int) (color.Palette, error) {
	histogram := q.CreateGreyscaleHistogram(img)
	ychan := make(chan []int)
	go multilevelThreshold(histogram, m, ychan, nil)
	T := <-ychan

	colours := make([]color.Color, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours, nil
}

func Colour(img image.Image) (color.Palette, error) {
	rhist, ghist, bhist, ahist := q.CreateRGBAHistogram(img)

	rchan := make(chan []int)
	gchan := make(chan []int)
	bchan := make(chan []int)
	achan := make(chan []int)

	var wg sync.WaitGroup
	wg.Add(4)

	go multilevelThreshold(rhist, 1, rchan, &wg)
	go multilevelThreshold(ghist, 1, gchan, &wg)
	go multilevelThreshold(bhist, 1, bchan, &wg)
	go multilevelThreshold(ahist, 1, achan, &wg)

	Tr := <-rchan
	Tg := <-gchan
	Tb := <-bchan
	Ta := <-achan

	wg.Wait()

	colours := make([]color.Color, len(Tr))
	for i := range colours {
		colours[i] = color.RGBA{uint8(Tr[i]), uint8(Tg[i]), uint8(Tb[i]), uint8(Ta[i])}
	}

	return colours, nil
}

func multilevelThreshold(hist q.Histogram, m int, c chan []int, wg *sync.WaitGroup) {
	xMax := 255
	xMin := 0

	// Calculate the initial threshold values
	T := make([]uint8, m+1)
	for i := 0; i <= m; i++ {
		T[i] = uint8(xMin + (i*(xMax-xMin))/m)
	}
	// Initialising the segment map
	segments := make(map[int]q.Histogram)
	for i := 1; i <= m; i++ {
		segments[i] = q.Histogram{}
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

func mean(h q.Histogram) int {
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

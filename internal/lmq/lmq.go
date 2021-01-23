package LMQ

import (
	"errors"
	q "github.com/nadav-rahimi/dominant-colour/internal/general"
	"image"
	"image/color"
	"sync"
)

type LMQ struct{}

func (lmq LMQ) Greyscale(img image.Image, m int) (color.Palette, error) {
	histogram := q.CreateGreyscaleHistogram(img)
	ychan := make(chan []int)
	go quantise(histogram, m, ychan, nil)
	T := <-ychan

	colours := make([]color.Color, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours, nil
}

func (lmq LMQ) Colour(img image.Image, m int) (color.Palette, error) {
	return nil, errors.New("LMQ does not support colour quantisation")
}

func quantise(hist q.Histogram, m int, c chan []int, wg *sync.WaitGroup) {
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

package lmq

import (
	"github.com/nadav-rahimi/dominant-colour/internal/quantisers"
	"image"
	"image/color"
)

const (
	xMax = 255
	xMin = 0
)

// Returns "m" greyscale colours to best recreate the colour palette of the original image
func QuantiseGreyscale(img image.Image, m int) color.Palette {
	// Create the histogram
	histogram := quantisers.CreateGreyscaleHistogram(img)

	// Calculate the initial threshold values
	T := make([]uint8, m+1)
	for i := 0; i <= m; i++ {
		T[i] = uint8(xMin + (i*(xMax-xMin))/m)
	}
	// Initialising the segment map
	segments := make(map[int]quantisers.LinearHistogram)
	for i := 1; i <= m; i++ {
		segments[i] = quantisers.LinearHistogram{}
	}
	// Initialising the averages for each segment
	averages := make(map[int]int)
	// Initialising the slice for the old threshold history
	oldT := make([]uint8, len(T))

	// Calculating the thresholds
	for {
		copy(oldT, T)

		// Segments the pixels of the image into thresholds based on histogram
		for k, v := range histogram {
			// Checks for k=0 since it cannot be checked in the loop
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

	// Convert the calculated averages to colour values
	colours := make([]color.Color, m)
	for i := range colours {
		colours[i] = color.Gray{uint8(averages[i+1])}
	}

	return colours
}

// Calculates the mean greyscale value in the histogram
func mean(h quantisers.LinearHistogram) int {
	sum, total := 0, 0

	for k, v := range h {
		sum += int(k) * v
		total += v
	}

	if total == 0 {
		return 0
	}
	return sum / total
}

// Checks if two uint8 slices are equal
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

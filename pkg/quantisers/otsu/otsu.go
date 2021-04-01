package otsu

import (
	"github.com/fiwippi/go-quantise/internal/quantisers"
	"image"
	"image/color"
)

const xMax = 256

// Returns one greyscale colour which best represents the threshold
// for splitting the image into black and white. Otsu only supports m = 1.
func QuantiseGreyscale(img image.Image) color.Palette {
	histogram := quantisers.CreateGreyscaleHistogram(img)
	threshold := calculateThreshold(histogram)
	return color.Palette{color.Gray{threshold}}
}

// Calculates the threshold for otsu for which to split colours on,
// all pixels which value below the threshold should be black and above
// the threshold should be white
func calculateThreshold(hist quantisers.LinearHistogram) uint8 {
	P := make([]int, xMax)
	S := make([]int, xMax)
	P[0], S[0] = 0, 0

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
			Sp := P[v] - P[u-1]
			Ss := S[v] - S[u-1]
			if P[v]-P[u-1] == 0 {
				H[u][v] = 0
				continue
			}

			H[u][v] = (Ss / Sp) * Ss
		}
	}

	var T uint8 = 0
	var maxVariation int = 0
	for t := uint8(0); t < xMax-1; t++ {
		variation := H[1][t] + H[t+1][xMax-1]
		if variation > maxVariation {
			maxVariation = variation
			T = t
		}
	}

	return T
}

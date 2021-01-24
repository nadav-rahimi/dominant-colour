package PNNLAB

import (
	"github.com/nadav-rahimi/dominant-colour/pkg/colours"
	"math"
)

// Node holding Colour data for PNN
type Node struct {
	Prev  *Node   // Pointer to the previous Node
	Next  *Node   // Pointer to the Next Node
	D     float64 // Merge cost value, indicating the increase in the MSE if the two classes are merged (this class and the one to the right)
	N     float64 // Number of pixels in the class
	Index int     // Index of the Node in the heap

	// Variables for greyscale quantisation
	C float64 // Mean grey level of the class
	T uint8   // Maximal grey value, also servers as threshold between the class and its neighbour class to the right

	// Variables for colour quantisation
	colours.RGB
	NN *Node // Pointers to the nearest neighbour
}

// Calculates the cost of merging two greyscale clusters,
// it represents the increase in MSE value caused by the merge
func LinearCost(a, b *Node) float64 {
	lhs := (a.N * b.N) / (a.N + b.N)
	rhs := math.Pow(math.Abs(a.C-b.C), 2)

	return lhs * rhs
}

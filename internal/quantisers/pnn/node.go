package pnn

import (
	_ "fmt"
	"github.com/fiwippi/go-quantise/pkg/colours"
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
	T uint8   // Maximal grey value, also serves as threshold between the class and its neighbour class to the right

	// Variables for colour quantisation
	colours.RGB         // RGB Values of the node
	A           float64 // Alpha Value of the node (Used for non-LAB PNN)
	NN          *Node   // Pointers to the nearest neighbour
	MergeCount  int     // The iteration where the node was last merged with another
	UpdateCount int     // The iteration where the MSE was last calculated for the node
}

// Squares a float64 number
func Sqr(a float64) float64 {
	return a * a
}

// Calculates the cost of merging two greyscale clusters,
// it represents the increase in MSE value caused by the merge
func LinearCost(a, b *Node) float64 {
	lhs := (a.N * b.N) / (a.N + b.N)
	rhs := Sqr(math.Abs(a.C - b.C))

	return lhs * rhs
}

// Calculates the cost of merging two colour clusters,
// it represents the increase in MSE value caused by the merge
func VectorCost(a, b *Node) float64 {
	lhs := (a.N * b.N) / (a.N + b.N)
	rhs := Sqr(b.A-a.A) + Sqr(b.R-a.R) + Sqr(b.G-a.G) + Sqr(b.B-a.B)

	return lhs * rhs
}

package pnn

import (
	_ "fmt"
	"math"
)

type Node struct {
	Prev  *Node   // Pointer to the previous Node
	Next  *Node   // Pointer to the Next Node
	C     float64 // Mean grey level of the class
	T     uint8   // Maximal grey value, also servers as threshold between the class and its neighbour class to the right
	D     float64 // Merge cost value, indicating the increase in the MSE if the two classes are merged (this class and the one to the right)
	N     float64 // Number of pixels in the class
	Index int     // Index of the Node in the heap

	// Extra variables for colour node
	A, R, G, B float64 // ARGB Value of the node
	NN         *Node   // Pointers to the nearest neighbour
}

func Sqr(a float64) float64 {
	return a * a
}

// Calculates the cost of merging two clusters, it represents the increase in MSE value caused by the merge
func LinearCost(a, b *Node) float64 {
	lhs := (a.N * b.N) / (a.N + b.N)
	rhs := math.Pow(math.Abs(a.C-b.C), 2)

	return lhs * rhs
}

func VectorCost(a, b *Node) float64 {
	lhs := (a.N * b.N) / (a.N + b.N)
	rhs := Sqr(b.A-a.A) + Sqr(b.R-a.R) + Sqr(b.G-a.G) + Sqr(b.B-a.B)

	return lhs * float64(rhs)
}

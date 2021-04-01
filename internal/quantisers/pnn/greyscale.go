package pnn

import (
	"container/heap"
	"github.com/fiwippi/go-quantise/internal/quantisers"
	"image"
	"image/color"
	"sort"
)

// Returns "m" greyscale colours to best recreate the colour palette of the original image
func QuantiseGreyscale(img image.Image, m int) color.Palette {
	hist := quantisers.CreateGreyscaleHistogram(img)
	T := calculateGreyscaleThresholds(hist, m)

	colours := make([]color.Color, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours
}

// Calculates the greyscale thresholds for the histogram
func calculateGreyscaleThresholds(hist quantisers.LinearHistogram, M int) []int {
	// Create linked list and heap for PNN
	S, H := initialiseGreyscaleStructures(hist)

	// m is set to H.Len() + 1 since the final element
	// in the list is left out from the heap so for total
	// number of elements, 1 needs to be added
	m := H.Len() + 1
	for m != M {
		sa := H.Front().(*Node)
		updateGreyscaleStructs(sa, sa.Next, H)
		m = m - 1
	}

	// Return the greyscale thresholds
	thresholds := make([]int, 0, M)
	for S != nil {
		thresholds = append(thresholds, int(S.C))
		S = S.Next
	}
	sort.Ints(thresholds)

	return thresholds
}

// Initialises the linked list and heap used by the PNN Algorithm to quantise the image
func initialiseGreyscaleStructures(hist quantisers.LinearHistogram) (*Node, *Heap) {
	// Initialise Heap
	h := make(Heap, 0)
	heap.Init(&h)

	// Initialise List
	var head *Node
	var currentNode *Node
	var previousNode *Node

	// Sort keys so linked list created
	// in order of increasing grey value
	keys := make([]int, 0)
	for k, _ := range hist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	previousNode = nil
	for i, k := range keys {
		// Create a new node
		currentNode = &Node{
			Prev:  previousNode,
			C:     float64(k),
			T:     uint8(k),
			D:     -1,
			N:     float64(hist[uint8(k)]),
			Index: -1,
		}

		if i == 0 {
			head = currentNode
		}

		if previousNode != nil {
			// Add the node to the list and calculate its cost
			previousNode.Next = currentNode
			previousNode.D = LinearCost(previousNode, currentNode)

			// Add the previous node to the heap
			heap.Push(&h, previousNode)
		}

		// Make the current node the next previous node
		previousNode = currentNode
	}

	return head, &h
}

// Reduces the size of the linked list to eventually achieve a quantised palette
func updateGreyscaleStructs(a, b *Node, h *Heap) {
	// Combine the data from B into A
	Nq := a.N + b.N
	Cq := (a.N*a.C + b.N*b.C) / Nq
	Tq := b.T

	a.N = Nq
	a.C = Cq
	a.T = Tq

	// If A is the penultimate element then its next must be set to nil
	// Otherwise remove B from the linked list
	if b.Next == nil {
		a.Next = nil
	} else {
		a.Next = b.Next
		b.Next.Prev = a
	}

	// Recalculate the MSE costs and update their locations in the heap with the new cost
	if a.Prev != nil {
		APrevCost := LinearCost(a.Prev, a)
		h.Update(a.Prev, APrevCost)
	}
	if a.Next != nil {
		ACost := LinearCost(a, a.Next)
		h.Update(a, ACost)
	}

	// Remove the second element from the heap if its in the heap
	if b.Index >= 0 && b.Index < h.Len() {
		_ = heap.Remove(h, b.Index)
	}

	// Remove element from heap if its at the end of the list
	if a.Next == nil {
		_ = heap.Remove(h, a.Index)
	}
}

package pnn

import (
	"container/heap"
	"github.com/nadav-rahimi/dominant-colour/pkg/colours"
	"image"
	"image/color"
	"math"
	"sort"
)

// Used as a variable for each pnn operation to determine what type of distance calculation to use,
// this variable is available for the whole scope of the pnn operation
type PNNMode uint8

// Which colour mode to use when calculating the distances between colours
const (
	RGB PNNMode = iota
	LAB
)

// Quantises a given into a palette of "m" colours to best represent it
func (mode PNNMode) QuantiseColour(img image.Image, M int) color.Palette {
	// Creates the histogram of the image
	hist := CreatePNNHistogram(img)

	// Make the linked list of nodes
	S, H := mode.initialiseColours(hist)

	m := H.Len() + 1
	count := 0
	for m != M {
		n := mode.recalculateNeighbours(H, count)
		mode.updateColourStructs(n, n.NN, H, count)

		m = m - 1
		count += 1
	}

	thresholds := make(color.Palette, 0, M)
	for S != nil {
		clr := color.RGBA{uint8(S.R), uint8(S.G), uint8(S.B), uint8(S.A)}
		thresholds = append(thresholds, clr)
		S = S.Next
	}

	return thresholds
}

// Recalculates nearest neighbours
func (mode PNNMode) recalculateNeighbours(H *Heap, count int) *Node {
	for {
		S := H.Front().(*Node)

		if S.UpdateCount >= S.MergeCount && S.UpdateCount >= S.NN.MergeCount {
			return S
		} else {
			mode.nearestNeighbour(S)
			heap.Fix(H, S.Index)
			S.UpdateCount = count
		}
	}
}

// Initialises the linked list and heap used by the PNN Algorithm to quantise the image
func (mode PNNMode) initialiseColours(hist Histogram) (*Node, *Heap) {
	// Initialise List
	var currentNode *Node
	var previousNode *Node

	keys := make([]int, 0)
	for k, _ := range hist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	//fmt.Println("Number of Colours:", len(keys))

	var head = hist[uint32(keys[0])]
	previousNode = nil
	for _, i := range keys {
		currentNode = hist[uint32(i)]
		currentNode.A /= currentNode.N
		currentNode.R /= currentNode.N
		currentNode.G /= currentNode.N
		currentNode.B /= currentNode.N

		currentNode.Prev = previousNode
		if previousNode != nil {
			previousNode.Next = currentNode
		}

		// Make the current node the next previous node
		previousNode = currentNode
	}

	// Make the heap
	h := make(Heap, 0)
	heap.Init(&h)

	// Initialise nearest neighbour for each node and build heap of nodes
	n := head
	for n != nil {
		mode.nearestNeighbour(n)
		if n.Next != nil {
			heap.Push(&h, n)
		}
		n = n.Next
	}

	return head, &h
}

// For a given node, finds the nearest neighbour, this is the node which has the smallest merge cost
func (mode PNNMode) nearestNeighbour(node *Node) {
	var err = math.MaxFloat64
	var nn *Node

	if mode == LAB {
		lab1 := node.LAB()
		tmp := node.Next
		for tmp != nil {
			lab2 := tmp.LAB()
			nerr := colours.LABDistance(lab1, lab2)
			if nerr < err {
				err = nerr
				nn = tmp
			}
			tmp = tmp.Next
		}
	} else { // Default to RGB if not LAB or any other mode
		tmp := node.Next
		for tmp != nil {
			nerr := VectorCost(node, tmp)
			if nerr < err {
				err = nerr
				nn = tmp
			}
			tmp = tmp.Next
		}
	}

	node.NN = nn
	node.D = err
}

// Reduces the size of the linked list to eventually achieve a quantised palette
func (mode PNNMode) updateColourStructs(a, b *Node, h *Heap, count int) {
	Nq := a.N + b.N
	a.A = (a.N*a.A + b.N*b.A) / Nq
	a.R = (a.N*a.R + b.N*b.R) / Nq
	a.G = (a.N*a.G + b.N*b.G) / Nq
	a.B = (a.N*a.B + b.N*b.B) / Nq
	a.N = Nq

	// Unchain the nearest neighbour bin
	if b.Next != nil {
		b.Next.Prev = b.Prev
	}
	if b.Prev != nil {
		//fmt.Printf("Removing %p %+v\n", b, b)
		b.Prev.Next = b.Next
	}

	// Remove the neighbour from the bin
	if b.Index >= 0 && b.Index < h.Len() {
		_ = heap.Remove(h, b.Index)
	}

	// Remove element from heap if its at the
	// end of the list and not already removed
	if a.Next == nil && a.Index != -1 {
		_ = heap.Remove(h, a.Index)
	}

	a.MergeCount = count + 1
	b.MergeCount = math.MaxInt32
}

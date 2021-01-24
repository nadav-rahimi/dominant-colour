package PNNLAB

import (
	"container/heap"
	"fmt"
	"github.com/nadav-rahimi/dominant-colour/pkg/colours"
	"image"
	"image/color"
	"math"
	"sort"
	"sync"
)

// TODO, more efficient way to calculate nearest neighbour

// Returns "m" rgba colours to best recreate the colour palette of the original image
func (pnn *PNNLAB) Colour(img image.Image, m int) (color.Palette, error) {
	hist := CreatePNNLABHistogram(img)
	ychan := make(chan color.Palette)
	go quantiseColour(hist, m, ychan, nil)
	colours := <-ychan

	return colours, nil
}

// Quantises a linear histogram with pnn
func quantiseColour(hist ColourHistogram, M int, c chan color.Palette, wg *sync.WaitGroup) {
	// Make the linked list of nodes
	S, H := initialiseColours(hist)

	fmt.Println("initialised")

	m := H.Len() + 1
	for m != M {
		fmt.Println(m)
		sa := H.Front().(*Node)
		updateColourStructs(sa, sa.NN, H)
		recalculateNeighbours(S, H)
		m = m - 1
	}

	thresholds := make(color.Palette, 0, M)
	for S != nil {
		c := color.RGBA{uint8(S.R), uint8(S.G), uint8(S.B), 255}
		thresholds = append(thresholds, c)
		S = S.Next
	}

	if wg != nil {
		wg.Done()
	}

	c <- thresholds
}

// Recalculates nearest neighbours for all nodes in the heap
func recalculateNeighbours(S *Node, H *Heap) {
	// Initialise nearest neighbour for each node and build heap of nodes
	n := S
	for n != nil {
		nearestNeighbour(n)
		heap.Fix(H, n.Index)
		n = n.Next
	}
}

// Initialises the linked list and heap used by the PNN Algorithm to quantise the image
func initialiseColours(hist ColourHistogram) (*Node, *Heap) {
	// Initialise List
	var currentNode *Node
	var previousNode *Node

	keys := make([]int, 0)
	for k, _ := range hist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var head = hist[uint32(keys[0])]
	previousNode = nil
	for _, i := range keys {
		currentNode = hist[uint32(i)]
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

	fmt.Println(len(keys))
	cnt := 1
	// Initialise nearest neighbour for each node and build heap of nodes
	n := head
	for n != nil {
		nearestNeighbour(n)
		cnt++
		if n.Next != nil {
			heap.Push(&h, n)
		}
		n = n.Next

	}

	return head, &h
}

// For a given node, finds the nearest neighbour, this is the node which has the smallest merge cost
func nearestNeighbour(node *Node) {
	var err = math.MaxFloat64
	var nn *Node

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

	node.NN = nn
	node.D = err
}

// Reduces the size of the linked list to eventually achieve a quantised palette
func updateColourStructs(a, b *Node, h *Heap) {
	Nq := a.N + b.N
	a.R = (a.N*a.R + b.N*b.R) / Nq
	a.G = (a.N*a.G + b.N*b.G) / Nq
	a.B = (a.N*a.B + b.N*b.B) / Nq
	a.N = Nq

	// Unchain the nearest neighbour bin
	if b.Next != nil {
		b.Next.Prev = b.Prev
	}
	if b.Prev != nil {
		b.Prev.Next = b.Next
	}

	// Remove the neighbour from the bin
	if b.Index >= 0 && b.Index < h.Len() {
		_ = heap.Remove(h, b.Index)
	}

	// Remove element from heap if its at the end of the list
	if a.Next == nil {
		_ = heap.Remove(h, a.Index)
	}
}

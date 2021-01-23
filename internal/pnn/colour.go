package PNN

import (
	"container/heap"
	"image"
	"image/color"
	"math"
	"sort"
	"sync"
)

func (pnn *PNN) Colour(img image.Image, m int) (color.Palette, error) {
	hist := CreatePNNHistogram(img)
	ychan := make(chan color.Palette)
	go quantiseColour(hist, m, ychan, nil)
	colours := <-ychan

	return colours, nil
}

func quantiseColour(hist ColourHistogram, M int, c chan color.Palette, wg *sync.WaitGroup) {
	// Make the linked list of nodes
	S, H := initialiseColours(hist)

	m := H.Len() + 1
	for m != M {
		sa := H.Front().(*Node)
		updateColourStructs(sa, sa.NN, H)
		recalculateNeighbours(S, H)
		m = m - 1
	}

	//fmt.Println("Done")

	thresholds := make(color.Palette, 0, M)
	for S != nil {
		c := color.RGBA{uint8(S.R), uint8(S.G), uint8(S.B), uint8(S.A)}
		//fmt.Println(c)
		thresholds = append(thresholds, c)
		S = S.Next
	}

	if wg != nil {
		wg.Done()
	}

	c <- thresholds
}

func recalculateNeighbours(S *Node, H *Heap) {
	// Initialise nearest neighbour for each node and build heap of nodes
	n := S
	for n != nil {
		nearestNeighbour(n)
		heap.Fix(H, n.Index)
		n = n.Next
	}
}

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
		nearestNeighbour(n)
		if n.Next != nil {
			heap.Push(&h, n)
		}
		n = n.Next
	}

	return head, &h
}

func nearestNeighbour(node *Node) {
	var err = math.MaxFloat64
	var nn *Node

	tmp := node.Next
	for tmp != nil {
		nerr := VectorCost(node, tmp)
		if nerr < err {
			err = nerr
			nn = tmp
		}
		tmp = tmp.Next
	}

	node.NN = nn
	node.D = err

}

func updateColourStructs(a, b *Node, h *Heap) {
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

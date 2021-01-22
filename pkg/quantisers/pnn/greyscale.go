package PNN

import (
	"container/heap"
	q "github.com/nadav-rahimi/dominant-colour/internal/quantisers"
	"github.com/nadav-rahimi/dominant-colour/internal/quantisers/pnn"
	"image"
	"image/color"
	"sort"
	"sync"
)

func Greyscale(img image.Image, m int) (color.Palette, error) {
	hist := q.CreateGreyscaleHistogram(img)
	ychan := make(chan []int)
	go threshold(hist, m, ychan, nil)
	T := <-ychan

	colours := make([]color.Color, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours, nil
}

func threshold(hist q.Histogram, M int, c chan []int, wg *sync.WaitGroup) {
	S, H := initialise(hist)

	m := H.Len() + 1
	for m != M {
		sa := H.Front().(*pnn.Node)
		updateDataStructs(sa, sa.Next, H)
		m = m - 1
	}

	thresholds := make([]int, 0, M)
	for S != nil {
		thresholds = append(thresholds, int(S.C))
		S = S.Next
	}
	sort.Ints(thresholds)

	if wg != nil {
		wg.Done()
	}

	c <- thresholds

}

func initialise(hist q.Histogram) (*pnn.Node, *pnn.Heap) {
	// Initialise Heap
	h := make(pnn.Heap, 0)
	heap.Init(&h)

	// Initialise List
	var head *pnn.Node
	var currentNode *pnn.Node
	var previousNode *pnn.Node

	keys := make([]int, 0)
	for k, _ := range hist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	previousNode = nil
	for i, k := range keys {
		// Create a new node
		currentNode = &pnn.Node{
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
			previousNode.D = pnn.LinearCost(previousNode, currentNode)

			// Add the previous node to the heap
			heap.Push(&h, previousNode)
		}

		// Make the current node the next previous node
		previousNode = currentNode
	}

	return head, &h
}

func updateDataStructs(a, b *pnn.Node, h *pnn.Heap) {
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
		APrevCost := pnn.LinearCost(a.Prev, a)
		h.Update(a.Prev, APrevCost)
	}
	if a.Next != nil {
		ACost := pnn.LinearCost(a, a.Next)
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

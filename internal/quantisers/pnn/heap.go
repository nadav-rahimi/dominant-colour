package pnn

// Adapted from https://golang.org/pkg/container/heap/

import "container/heap"

type Heap []*Node

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	// We want Pop to give us the lowest priority so we use less than here.
	return h[i].D < h[j].D
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *Heap) Push(x interface{}) {
	n := len(*h)
	item := x.(*Node)
	item.Index = n
	*h = append(*h, item)
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

// update modifies the priority (the merge cost value) of an Item in the queue.
func (h *Heap) Update(node *Node, d float64) {
	node.D = d
	heap.Fix(h, node.Index)
}

package pnn

import "container/heap"

// Implements a Heap of Nodes used for PNN
type Heap []*Node

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	// We want Pop to give us the node with lowest cost so we use "<" here.
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

func (h *Heap) Front() interface{} {
	return (*h)[0]
}

func (h *Heap) Update(node *Node, d float64) {
	node.D = d
	heap.Fix(h, node.Index)
}

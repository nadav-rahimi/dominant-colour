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
	bin := q.CreateGreyscaleBin(img)
	ychan := make(chan []int)
	go threshold(bin, m, ychan, nil)
	T := <-ychan

	colours := make([]color.Color, len(T))
	for i := range colours {
		colours[i] = color.Gray{uint8(T[i])}
	}

	return colours, nil
}

func threshold(bin q.Bin, M int, c chan []int, wg *sync.WaitGroup) {
	_, H := initialise(bin)

	//for S, _ = initialise(bin); S != nil; S = S.Next {
	//	fmt.Printf("%+v\n", S)
	//}

	m := len(bin) // TODO set this to length of S or lengh of heap if quicker?
	for m != M {
		sa := heap.Pop(H).(*pnn.Node)
		sb := sa.Next
		heap.Push(H, sa) // Need to push it back on because we do not want to change the list yet

		//fmt.Printf("%+v, %+v\n", sa, sb)

		merge(sa, sb)
		updateDataStructs(sa, sb, H)
		m = m - 1
	}

	thresholds := make([]int, 0, M)
	for H.Len() > 0 {
		node := heap.Pop(H).(*pnn.Node)
		thresholds = append(thresholds, int(node.T))
	}
	sort.Ints(thresholds)

	if wg != nil {
		wg.Done()
	}

	c <- thresholds

}

func initialise(bin q.Bin) (*pnn.Node, *pnn.Heap) {
	// Initialise Heap
	h := make(pnn.Heap, 0)
	heap.Init(&h)

	// Initialise List
	var head *pnn.Node
	var currentNode *pnn.Node
	var previousNode *pnn.Node

	// Initialise the head separately since its previous node will be nil
	head = &pnn.Node{
		Prev: nil,
		C:    0,
		T:    0,
		D:    -1,
		N:    float64(bin[0]),
	}
	heap.Push(&h, head)
	previousNode = head

	//fmt.Println(len(bin))
	for k, v := range bin {
		if k == 0 {
			continue
		}

		// Create a new node
		currentNode = &pnn.Node{
			Prev: previousNode,
			C:    float64(k),
			T:    float64(k),
			D:    -1,
			N:    float64(v),
		}

		// Add the node to the list and calculate its cost
		previousNode.Next = currentNode
		previousNode.D = pnn.Cost(previousNode, currentNode)

		// Add the current node to the heap
		heap.Push(&h, currentNode)

		// Make the current node the next previous node
		previousNode = currentNode
	}

	// TODO TEST FIX
	head.Prev = previousNode
	previousNode.Next = head

	return head, &h
}

func merge(a, b *pnn.Node) {
	Nq := a.N + b.N
	Cq := (a.N*a.C + b.N*b.C) / Nq
	Tq := b.T

	a.N = Nq
	a.C = Cq
	a.T = Tq
}

func updateDataStructs(a, b *pnn.Node, h *pnn.Heap) {
	// Remove the second element from the linked list
	a.Next = b.Next
	b.Next.Prev = a

	// Recalculate the MSE costs
	a.Prev.D = pnn.Cost(a.Prev, a)
	a.D = pnn.Cost(a, a.Next)

	// Remove the second element from the heap
	//fmt.Println(b.Index)
	_ = heap.Remove(h, b.Index)

	// Update the locations in the heap of the elements with the new cost
	h.Update(a.Prev, a.Prev.D)
	h.Update(a, a.D)
}

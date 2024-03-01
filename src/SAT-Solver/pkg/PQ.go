package pkg

import (
	"container/heap"
)

// Item represents an item in the priority queue.
type Item struct {
	value    int // The value of the item
	priority int // The priority of the item
	index    int // The index of the item in the heap
}

// PriorityQueue implements a priority queue based on a min-heap.
type PriorityQueue []*Item

// Len returns the length of the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares the priorities of two items.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

// Swap swaps two items in the priority queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Update modifies the priority and value of an item in the priority queue.
func (pq *PriorityQueue) Update(item *Item, value int, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

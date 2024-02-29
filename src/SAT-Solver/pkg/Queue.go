package main

import "fmt"

// Queue represents a queue data structure
type Queue []int

// Enqueue adds an element to the end of the queue
func (q *Queue) Enqueue(item int) {
	*q = append(*q, item)
}

// Dequeue removes and returns the first element from the queue
func (q *Queue) Dequeue() (int, error) {
	if len(*q) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, nil
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

func main() {
	q := make(Queue, 0)

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	for !q.IsEmpty() {
		item, err := q.Dequeue()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Dequeued:", item)
	}
}

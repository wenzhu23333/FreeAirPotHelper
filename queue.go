package main

import (
	"fmt"
	"sync"
)

type Queue struct {
	items    []string
	capacity int
	lock     sync.Mutex
}

func NewQueue(capacity int) *Queue {
	q := &Queue{
		items:    make([]string, 0, capacity),
		capacity: capacity,
	}
	return q
}

func (q *Queue) Enqueue(item string) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.items) == q.capacity {
		q.items = q.items[1:]
	}
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (string, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return "", fmt.Errorf("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *Queue) Head() (string, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) == 0 {
		return "", fmt.Errorf("queue is empty")
	}
	return q.items[0], nil
}

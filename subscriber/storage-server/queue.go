package main

import (
	"container/list"
)

type Queue struct {
	queue *list.List
}

func NewQueue() *Queue {
	q := &Queue{queue: list.New()}
	return q
}

func (q *Queue) Put(v interface{}) *list.Element {
	return q.queue.PushBack(v)
}

func (q *Queue) Get() *list.Element {
	item := q.queue.Front()
	q.queue.Remove(item)
	return item
}

func (q *Queue) Len() int {
	return q.queue.Len()
}

func (q *Queue) Empty() bool {
	return q.Len() == 0
}

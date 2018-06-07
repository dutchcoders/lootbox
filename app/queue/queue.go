package queue

import "net/url"

var (
	queue = NewUrlStack()
)

func Push(v *url.URL) {
	queue.Push(v)
}

func Pop() *url.URL {
	return queue.Pop()
}

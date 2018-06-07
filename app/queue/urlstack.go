package queue

import (
	"net/url"
	"time"

	"github.com/ReneKroon/ttlcache"
)

type UrlStack struct {
	*Stack
	cache *ttlcache.Cache
}

func NewUrlStack() *UrlStack {
	q := NewStack()

	cache := ttlcache.NewCache()
	cache.SetTTL(time.Duration(60 * time.Hour))

	return &UrlStack{
		Stack: q,
		cache: cache,
	}
}

func (q *UrlStack) Push(v *url.URL) {
	if _, exists := q.cache.Get(v.String()); exists {
		return
	}

	q.cache.Set(v.String(), v.String())

	q.Stack.Push(v)
}

func (q *UrlStack) Pop() *url.URL {
	v := q.Stack.Pop()
	if v == nil {
		return nil

	}
	return v.(*url.URL)
}

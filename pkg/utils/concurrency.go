package utils

import (
	"sync"
)

type Broadcaster struct {
	mu        sync.Mutex
	listeners map[chan struct{}]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		listeners: make(map[chan struct{}]struct{}),
	}
}

func (b *Broadcaster) Subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	b.mu.Lock()
	b.listeners[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) Unsubscribe(ch chan struct{}) {
	b.mu.Lock()
	if _, ok := b.listeners[ch]; ok {
		delete(b.listeners, ch)
		close(ch)
	}
	b.mu.Unlock()
}

func (b *Broadcaster) Broadcast() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.listeners {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

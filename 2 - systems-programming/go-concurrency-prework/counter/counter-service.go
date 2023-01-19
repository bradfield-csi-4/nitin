package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type CounterService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently from multiple goroutines without any
	// additional synchronization.
	getNext() uint64
}

type UnsafeCounter struct {
	value uint64
}

func (c *UnsafeCounter) getNext() uint64 {
	c.value++
	return c.value
}

type AtomicCounter struct {
	value uint64
}

func (c *AtomicCounter) getNext() uint64 {
	return atomic.AddUint64(&c.value, 1)
}

type LockedCounter struct {
	value uint64
	mu    sync.Mutex
}

func (c *LockedCounter) getNext() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

type ChannelCounter struct {
	Value uint64
}

var requests = make(chan struct{})
var responses = make(chan uint64)

func (c *ChannelCounter) getNext() uint64 {
	requests <- struct{}{}
	return <-responses
}

func counter() {
	var counter uint64
	for {
		<-requests
		counter++
		responses <- counter
	}
}

func init() {
	go counter()
}

func main() {
	counter := ChannelCounter{}
	num := 10000

	wg := new(sync.WaitGroup)
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			counter.getNext()
		}(wg)
	}

	wg.Wait()

	fmt.Println(counter.getNext())
}
